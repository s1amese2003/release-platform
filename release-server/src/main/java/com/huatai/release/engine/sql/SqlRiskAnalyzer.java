package com.huatai.release.engine.sql;

import com.huatai.release.engine.sql.rule.AlterRiskRule;
import com.huatai.release.engine.sql.rule.CreateIfNotExistsRule;
import com.huatai.release.engine.sql.rule.CrossDatabaseRule;
import com.huatai.release.engine.sql.rule.DmlWhereRule;
import com.huatai.release.engine.sql.rule.ForbiddenDdlRule;
import com.huatai.release.engine.sql.rule.PermissionRule;
import com.huatai.release.engine.sql.rule.SystemSchemaRule;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.dto.SqlAuditReport;
import com.huatai.release.model.enums.RiskLevel;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

@Component
public class SqlRiskAnalyzer {

    private final List<SqlRule> rules = List.of(
            new ForbiddenDdlRule(),
            new DmlWhereRule(),
            new SystemSchemaRule(),
            new AlterRiskRule(),
            new PermissionRule(),
            new CrossDatabaseRule(),
            new CreateIfNotExistsRule()
    );

    public SqlAuditReport analyze(String sqlContent, String version, String businessDb) {
        List<SqlBlock> blocks = splitSql(sqlContent);
        List<SqlAuditIssue> issues = new ArrayList<>();

        for (SqlBlock block : blocks) {
            SqlRuleContext context = new SqlRuleContext(block.sql(), block.line(), businessDb);
            for (SqlRule rule : rules) {
                rule.apply(context).ifPresent(issues::add);
            }
        }

        Set<Integer> riskyLines = new HashSet<>();
        issues.forEach(issue -> riskyLines.add(issue.line()));

        Map<String, Long> summary = new LinkedHashMap<>();
        summary.put("P0", issues.stream().filter(i -> i.level() == RiskLevel.P0).count());
        summary.put("P1", issues.stream().filter(i -> i.level() == RiskLevel.P1).count());
        summary.put("P2", issues.stream().filter(i -> i.level() == RiskLevel.P2).count());
        summary.put("SAFE", (long) blocks.size() - riskyLines.size());

        return new SqlAuditReport(version, blocks.size(), issues, summary);
    }

    private List<SqlBlock> splitSql(String sqlContent) {
        List<SqlBlock> result = new ArrayList<>();
        if (sqlContent == null || sqlContent.isBlank()) {
            return result;
        }

        StringBuilder sb = new StringBuilder();
        int line = 1;
        int beginLine = 1;
        for (char ch : sqlContent.toCharArray()) {
            if (sb.length() == 0 && !Character.isWhitespace(ch)) {
                beginLine = line;
            }
            sb.append(ch);

            if (ch == ';') {
                String sql = sb.toString().trim();
                if (!sql.isBlank()) {
                    result.add(new SqlBlock(sql, beginLine));
                }
                sb.setLength(0);
            }

            if (ch == '\n') {
                line++;
            }
        }

        String tail = sb.toString().trim();
        if (!tail.isBlank()) {
            result.add(new SqlBlock(tail, beginLine));
        }
        return result;
    }
}
