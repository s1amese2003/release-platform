# 椤圭洰涓婄嚎骞冲彴

## 鐩綍

- `release-server`锛歋pring Boot 3 鍚庣
- `release-web`锛歏ue3 + Element Plus 鍓嶇
- `docs`锛氳鍒欎笌鍩虹嚎鏂囨。

## 鍚庣鑳藉姏

- 涓婁紶 JAR/ZIP 骞惰В鏋愶細
  - `upgrade/{yyyyMMdd}` 鐗堟湰璇嗗埆
  - 鏈€鏂?`upgrade.sql` 涓?`config.txt` 鎻愬彇
  - `bootstrap*.yml` / `application*.yml` 瑙ｆ瀽
  - `BOOT-INF/lib` 渚濊禆鎻愬彇
  - `MANIFEST.MF` 鏋勫缓淇℃伅鎻愬彇
- SQL 瀹¤锛歅0/P1/P2 瑙勫垯寮曟搸
- 閰嶇疆姣斿锛氫笌鐜鍩虹嚎閫愰」 diff
- 渚濊禆鎵弿锛氱増鏈彉鏇?+ 甯歌楂樺嵄渚濊禆鎻愮ず
- 瀹℃壒鍔ㄤ綔锛欰PPROVE / REJECT / DEPLOY

## API

- `POST /api/releases/analyze`
- `GET /api/releases/{requestNo}`
- `POST /api/releases/{requestNo}/approval`
- `GET /api/baselines/{env}/{appName}`
- `POST /api/baselines/{env}/{appName}`
- `GET /api/system/health`

## 鍚姩

### 鍚庣

```bash
cd release-server
mvn spring-boot:run
```

### 鍓嶇

```bash
cd release-web
npm install
npm run dev
```

榛樿鑱旇皟鍦板潃锛氬墠绔?`5173` 浠ｇ悊鍒板悗绔?`8080`銆?
