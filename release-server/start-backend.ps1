$env:JAVA_HOME='C:\Program Files\Eclipse Adoptium\jdk-17.0.18.8-hotspot'
$env:PATH="$env:JAVA_HOME\bin;E:\release-platform\tools\apache-maven-3.9.12\bin;$env:PATH"
Set-Location 'E:\release-platform\release-server'
mvn spring-boot:run -DskipTests
