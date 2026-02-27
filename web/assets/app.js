const state = {
  releases: [],
  selectedReleaseId: null,
};

const dom = {
  healthBadge: document.getElementById("healthBadge"),
  createForm: document.getElementById("createForm"),
  versionInput: document.getElementById("versionInput"),
  releaseTableBody: document.querySelector("#releaseTable tbody"),
  selectedReleaseText: document.getElementById("selectedReleaseText"),
  riskStat: document.getElementById("riskStat"),
  configStat: document.getElementById("configStat"),
  resultStat: document.getElementById("resultStat"),
  sqlIssueBody: document.querySelector("#sqlIssueTable tbody"),
  configDiffBody: document.querySelector("#configDiffTable tbody"),
  logBox: document.getElementById("logBox"),
  artifactFile: document.getElementById("artifactFile"),
  uploadBtn: document.getElementById("uploadBtn"),
  scanBtn: document.getElementById("scanBtn"),
  approveBtn: document.getElementById("approveBtn"),
  deployBtn: document.getElementById("deployBtn"),
};

function nowVersion() {
  const d = new Date();
  const pad = (n) => `${n}`.padStart(2, "0");
  return `web-${d.getFullYear()}${pad(d.getMonth() + 1)}${pad(d.getDate())}-${pad(d.getHours())}${pad(d.getMinutes())}${pad(d.getSeconds())}`;
}

dom.versionInput.value = nowVersion();

function appendLog(message, isError = false) {
  const line = `[${new Date().toLocaleTimeString()}] ${isError ? "ERROR" : "INFO"} ${message}`;
  dom.logBox.textContent += `${line}\n`;
  dom.logBox.scrollTop = dom.logBox.scrollHeight;
}

async function api(path, options = {}) {
  const response = await fetch(path, options);
  const raw = await response.text();
  const data = raw ? JSON.parse(raw) : {};
  if (!response.ok) {
    throw new Error(data.error || `HTTP ${response.status}`);
  }
  return data;
}

function badge(text) {
  const lower = String(text || "").toLowerCase();
  let cls = "badge-pending";
  if (lower.includes("pass") || lower.includes("success") || lower.includes("approved")) cls = "badge-pass";
  if (lower.includes("warn") || lower.includes("running") || lower.includes("pending")) cls = "badge-warn";
  if (lower.includes("block") || lower.includes("failed") || lower.includes("reject")) cls = "badge-block";
  return `<span class="badge ${cls}">${text || "-"}</span>`;
}

async function checkHealth() {
  try {
    const result = await api("/healthz");
    dom.healthBadge.innerHTML = `API ${badge(result.status)}`;
  } catch (err) {
    dom.healthBadge.innerHTML = `API ${badge("DOWN")}`;
    appendLog(`Health check failed: ${err.message}`, true);
  }
}

function renderReleaseTable() {
  dom.releaseTableBody.innerHTML = "";
  for (const release of state.releases) {
    const tr = document.createElement("tr");
    if (release.id === state.selectedReleaseId) tr.classList.add("active");
    tr.innerHTML = `
      <td>${release.id}</td>
      <td>${release.Application?.name || "-"}</td>
      <td>${release.version}</td>
      <td>${badge(release.scan_status)}</td>
      <td>${badge(release.admission_result)}</td>
      <td>${badge(release.approval_status)}</td>
      <td>${badge(release.deploy_status)}</td>
    `;
    tr.addEventListener("click", () => selectRelease(release.id));
    dom.releaseTableBody.appendChild(tr);
  }
}

async function loadReleases(keepSelection = true) {
  const prev = keepSelection ? state.selectedReleaseId : null;
  const data = await api("/api/releases");
  state.releases = data.items || [];
  if (!prev && state.releases.length > 0) state.selectedReleaseId = state.releases[0].id;
  if (prev && state.releases.some((x) => x.id === prev)) state.selectedReleaseId = prev;
  renderReleaseTable();
  if (state.selectedReleaseId) await loadReport(state.selectedReleaseId);
}

function truncate(text, max = 90) {
  if (!text) return "";
  const clean = String(text).replace(/\s+/g, " ");
  return clean.length > max ? `${clean.slice(0, max)}...` : clean;
}

function renderReport(ticket) {
  dom.selectedReleaseText.textContent = `当前选择: #${ticket.id} | ${ticket.Application?.name || "-"} | ${ticket.version}`;

  const report = ticket.ScanReport || {};
  dom.riskStat.innerHTML = `HIGH ${report.high_count || 0} / MEDIUM ${report.medium_count || 0} / LOW ${report.low_count || 0}`;
  dom.configStat.innerHTML = `新增 ${report.config_added || 0} / 删除 ${report.config_deleted || 0} / 修改 ${report.config_modified || 0}`;
  dom.resultStat.innerHTML = `${badge(ticket.admission_result)} ${report.block_reason ? ` ${report.block_reason}` : ""}`;

  dom.sqlIssueBody.innerHTML = "";
  const issues = report.sql_issues || [];
  if (issues.length === 0) {
    dom.sqlIssueBody.innerHTML = '<tr><td colspan="4">无命中</td></tr>';
  } else {
    for (const issue of issues) {
      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td>${issue.sequence}</td>
        <td>${badge(issue.risk_level)}</td>
        <td>${issue.rule_name}</td>
        <td title="${String(issue.snippet || "").replaceAll('"', "'")}">${truncate(issue.snippet, 100)}</td>
      `;
      dom.sqlIssueBody.appendChild(tr);
    }
  }

  dom.configDiffBody.innerHTML = "";
  const diffs = report.config_diffs || [];
  if (diffs.length === 0) {
    dom.configDiffBody.innerHTML = '<tr><td colspan="5">无差异</td></tr>';
  } else {
    for (const diff of diffs) {
      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td>${diff.config_key}</td>
        <td>${truncate(diff.online_value, 40)}</td>
        <td>${truncate(diff.candidate_value, 40)}</td>
        <td>${badge(diff.diff_type)}</td>
        <td>${diff.is_critical ? badge("CRITICAL") : badge("NORMAL")}</td>
      `;
      dom.configDiffBody.appendChild(tr);
    }
  }
}

async function loadReport(releaseId) {
  const ticket = await api(`/api/releases/${releaseId}/report`);
  renderReport(ticket);
}

async function selectRelease(releaseId) {
  state.selectedReleaseId = releaseId;
  renderReleaseTable();
  await loadReport(releaseId);
}

function getSelectedReleaseId() {
  if (!state.selectedReleaseId) throw new Error("请先在列表里选择发布单");
  return state.selectedReleaseId;
}

dom.createForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    const formData = new FormData(dom.createForm);
    const payload = Object.fromEntries(formData.entries());
    const data = await api("/api/releases", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    state.selectedReleaseId = data.id;
    appendLog(`创建发布单成功: #${data.id}`);
    dom.versionInput.value = nowVersion();
    await loadReleases();
  } catch (err) {
    appendLog(`创建发布单失败: ${err.message}`, true);
  }
});

dom.uploadBtn.addEventListener("click", async () => {
  try {
    const releaseId = getSelectedReleaseId();
    if (!dom.artifactFile.files[0]) throw new Error("请选择 zip/jar/war 文件");
    const form = new FormData();
    form.append("artifact", dom.artifactFile.files[0]);
    await api(`/api/releases/${releaseId}/artifact`, { method: "POST", body: form });
    appendLog(`上传制品成功: release #${releaseId}`);
    await loadReleases();
  } catch (err) {
    appendLog(`上传失败: ${err.message}`, true);
  }
});

dom.scanBtn.addEventListener("click", async () => {
  try {
    const releaseId = getSelectedReleaseId();
    await api(`/api/releases/${releaseId}/scan`, {
      method: "POST",
      headers: { "X-User": "web-user" },
    });
    appendLog(`触发扫描: release #${releaseId}`);
    await loadReleases();
  } catch (err) {
    appendLog(`扫描失败: ${err.message}`, true);
  }
});

dom.approveBtn.addEventListener("click", async () => {
  try {
    const releaseId = getSelectedReleaseId();
    await api(`/api/releases/${releaseId}/approve`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Role": "owner",
      },
      body: JSON.stringify({
        approver: "web-approver",
        action: "approve",
        comment: "approved from web console",
        level: 2,
      }),
    });
    appendLog(`审批通过: release #${releaseId}`);
    await loadReleases();
  } catch (err) {
    appendLog(`审批失败: ${err.message}`, true);
  }
});

dom.deployBtn.addEventListener("click", async () => {
  try {
    const releaseId = getSelectedReleaseId();
    await api(`/api/releases/${releaseId}/deploy`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Role": "ops",
      },
      body: JSON.stringify({ operator: "web-ops" }),
    });
    appendLog(`部署完成: release #${releaseId}`);
    await loadReleases();
  } catch (err) {
    appendLog(`部署失败: ${err.message}`, true);
  }
});

(async function boot() {
  appendLog("页面初始化完成");
  await checkHealth();
  await loadReleases(false);
  setInterval(async () => {
    try {
      await loadReleases();
    } catch (err) {
      appendLog(`刷新列表失败: ${err.message}`, true);
    }
  }, 8000);
})();
