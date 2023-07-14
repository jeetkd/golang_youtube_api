# golang_youtube_api
golang으로 youtube에서 제공하는 API를 사용하여 유튜브 공개 플레이 리스트 목록의 제목과 아디를 가져오는 코드

### ⚙개발 환경
***
- `go 1.18`
- **IDE** : Goland 2022.1
- **External Package** : "google.golang.org/api/option", "google.golang.org/api/youtube/v3"
- **OS** : Windows10

### ✔사전 준비
***
- **자세한 설명 사이트 참조** : https://healer4-13.tistory.com/12
1. 구글 클라우드 플랫폼에서 프로젝트 생성
2. YouTube Data API v3 사용
3. API 키 발급

### 🗂외부 패키지 다운
***

    go get google.golang.org/api/youtube/v3

### 📃참고 사이트
***
- https://developers.google.com/youtube/v3/docs/playlists/list?hl=ko
- https://github.com/youtube/api-samples/tree/master/go
- https://developers.google.com/youtube/v3/code_samples/go?hl=ko

### 🔑playlists ID(플레이 리스트 아이디)
***
- **자세한 설명 사이트 참조** : https://healer4-13.tistory.com/13
1. 원하는 유튜브 검색
2. **"재생목록"** 클릭
3. 원하는 **"재생목록 보기"** 클릭
4. 위에 url창 확인 : https://www.youtube.com/playlist?list=PLp-ofiyo_L4kUdOxLd0Jn1-G9neeA4ZLX
5. **PLp-ofiyo_L4kUdOxLd0Jn1-G9neeA4ZLX** 복사 후 playlistID flag에 넣어줌
