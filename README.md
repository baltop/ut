아래는 config.yaml 예시입니다:

yaml

- url: "http://localhost:8080/test"
- interval_sec: 5
- data_format:
-   temperature: "float"
-   humidity: "int"
-   status: "string"
-   alert: "bool"

실행 방법:
위 config.yaml 파일을 저장합니다.

Go 코드를 main.go로 저장합니다.

터미널에서 아래 명령어를 실행합니다:

bash
복사
편집
go run main.go config.yaml
이 프로그램은 설정된 interval_sec 간격으로 난수 기반 JSON을 생성해 지정된 url로 전송합니다. 테스트용 API 서버를 미리 실행시켜야 정상 동작합니다.
