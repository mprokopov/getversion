# getversion
small go utility for CICD sysyems fetching project version from gradle.properties or package.json


## Usage
in `gradle.properties` you need to have 

`version=1.2.3`

execute

`getVersion -source=gradle` should return 1.2.3

in `project.json` you need to have 

`"version": "1.2.3"`


`getVersion -source=node` should returh 1.2.3


in Jenkinsfile you can call it like 

`def dockerTag = sh(script: 'getVersion -source=node', returnStdout: true) + '.' +BUILD_ID`
