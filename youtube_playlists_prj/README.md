# 🔎프로젝트 개요
유튜브 게시물을 업로드한 사용자가 게시물을 비공개 또는 삭제하면 내 재생목록으로 가져온 게시물도 이용할 수 없고, 알려 주지도 않는다.
그래서 youtube api를 사용하여 내 재생목록이 업데이트 되었는지 확인 후 알려주는 프로그램을 만들기로 하였다.

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
- **자세한 설명 사이트 참조** : https://healer4-13.tistory.com/16
1. 구글 클라우드 플랫폼에서 프로젝트 생성
2. YouTube Data API v3 사용
3. OAuth 동의화면 설정
4. OAuth 클라이언트 ID 발급
5. 클라이언트 보안 비밀번호 json으로 다운로드
6. 유튜브 계정에서 테스트 계정에 대한 권한 추가

### ☑️코드 실행 전 체크해야 될 부분
***
- **main.go**
1. func main() : playlistID = ""             //재생목록 ID 넣어주세요
2. func (y *youtubeinfolists) SendMail() :
   2.1. username := ""                       // 이메일을 보낼 구글 계정 입력
   2.2. passwd := ""                         // 구글 앱 비밀번호
   2.3. to := []string{""}                   //받는 사람
- **oauth2.go**
1. func GoogleLoginAuth(url string) :
   1.1. var email = ""                       // 로그인 이메일
   1.2. var password = ""                    // 이메일 비밀번호
- **_client.json 다운로드**
  
### 🗂외부 패키지 다운
***

    go get "google.golang.org/api/youtube/v3"
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

### ▶프로젝트 동작 구조 요약
***
- **Youtube 비공개 재생목록 정보를 가지고 오는 인증절차**
1. 사용자 인증을 위한 url 받음
2. 받은 url로 이동 후 구글 로그인(자동화)을 하여 사용자 인증 후 코드를 얻음
3. 인증 코드를 토큰으로 전환
- **유튜브 비공개 재생목록을 비교 하기 위한 절차**
1. 얻은 토큰으로 비공개 재생목록의 title(제목)과 id를 가져옴
2. 처음 실행 할 경우 : 가져온 재생목록의 title과 id를 저장하기 위해 WriteFile 메소드를 실행하여 .txt 파일을 만듬
3. 로컬(WriteFile메소드로 저장한 파일)에 있는 재생목록의 title과 id를 가져오기 위해 ReadFile 메소드 실행
4. 1번에서 얻은 재생목록과 3번에서 가져온 재생목록을 비교하여 삭제 또는 추가된 재생목록이 있는지 CheckPlaylists 메소드로 확인
5. 삭제 또는 추가할 리스트들을 로컬에 있는 파일에 UpdatePlaylists 메소드로 업데이트함(삭제된 재생목록이 있을시 메일로 알림)
