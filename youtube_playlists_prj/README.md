# golang_youtube_api
golang으로 youtube에서 제공하는 OAuth 클라이언트 ID를 사용하여 유튜브 비공개 플레이 리스트 목록의 제목과 아디를 가져오는 코드

### ⚙개발 환경
***
- `go1.21.0`
- **IDE** : Goland 2022.1
- **External Package**
  1. "google.golang.org/api/option"
  2. "google.golang.org/api/youtube/v3"
  3. "golang.org/x/oauth2"
  4. "github.com/Davincible/chromedp-undetected"
  5. "github.com/chromedp/chromedp"
- **OS** : Windows10

### ✔사전 준비
***
- **자세한 설명 사이트 참조** : https://healer4-13.tistory.com/14
1. 구글 클라우드 플랫폼에서 프로젝트 생성
2. YouTube Data API v3 사용
3. OAuth 동의화면 설정
4. OAuth 클라이언트 ID 발급
5. 클라이언트 보안 비밀번호 json으로 다운로드

### 🗂외부 패키지 다운
***

    go get google.golang.org/api/youtube/v3
    go get "golang.org/x/oauth2"
    go get "github.com/Davincible/chromedp-undetected"
    go get "github.com/chromedp/chromedp"

### 📃참고 사이트
***
- https://developers.google.com/youtube/v3/docs/playlists/list?hl=ko
- https://github.com/youtube/api-samples/tree/master/go
- https://developers.google.com/youtube/v3/code_samples/go?hl=ko
- https://www.joinc.co.kr/w/man/12/oAuth2/Google

### 🔑playlists ID(플레이 리스트 아이디)
***
- **자세한 설명 사이트 참조** : https://healer4-13.tistory.com/13
1. 원하는 유튜브 검색
2. **"재생목록"** 클릭
3. 원하는 **"재생목록 보기"** 클릭
4. 위에 url창 확인 : https://www.youtube.com/playlist?list=PLp-ofiyo_L4kUdOxLd0Jn1-G9neeA4ZLX
5. **PLp-ofiyo_L4kUdOxLd0Jn1-G9neeA4ZLX** 복사 후 playlistID에 넣어줌

### ▶실행 방법
***
- **자세한 설명 사이트 참조** : https://healer4-13.tistory.com/15
1. 사용자 인증을 위한 url 얻음
2. url을 통해서 사용자 인증 후 인증 코드 반환
3. 인증코드(토큰) 입력
