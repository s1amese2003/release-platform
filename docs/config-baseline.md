# 閰嶇疆鍩虹嚎璇存槑

## 鍩虹嚎妯″瀷

鍩虹嚎鎸?`env + appName` 缁村害瀛樺偍锛屽€间负鎵佸钩鍖?`key-value`锛?
- dev
- sit
- uat
- prod

褰撳墠鐗堟湰榛樿浣跨敤鍐呭瓨瀛樺偍锛坄BaselineMemoryStore`锛夛紝鍚庣画鍙浛鎹负 MySQL `env_baseline` 琛ㄣ€?
## 姣斿缁村害

- 鏁版嵁搴撹繛鎺ワ細`spring.datasource.*`
- Redis锛歚spring.redis.*`
- Nacos锛歚nacos.*`
- Profile锛歚spring.profiles.active`
- 绔彛锛歚server.port`
- 璋冭瘯/鏂囨。寮€鍏筹細`swagger.enabled`銆乣springdoc.api-docs.enabled`銆乣debug.*`
- 鏁忔劅椤规壂鎻忥細`password|secret|token|*.key`

## 鐢熶骇鐜寮虹害鏉?
- `spring.profiles.active` 蹇呴』涓?`prod`锛圥0锛?- 浠ヤ笅寮€鍏充笉鍙负 `true`锛圥1锛?  - `swagger.enabled`
  - `springdoc.api-docs.enabled`
  - `debug.iceCmdb.systemNames`
  - `debug.workflowMaintenace.adjustParticipant`

## 椤甸潰鎿嶄綔

鍓嶇 `鐜鍩虹嚎` 椤甸潰鏀寔锛?
1. 鎸夌幆澧冨拰搴旂敤鍔犺浇鍩虹嚎銆?2. 鏂囨湰鏍煎紡 `key=value` 缂栬緫銆?3. 涓€閿繚瀛樺洖鍚庣銆?
