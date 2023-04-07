package handler

// 맵에 유저의 정보를 담아줄 함수
func UserSet1(user models.GithubUserInfo, m map[string]int) {
	m[user.Login] = 1
}

// 맵에 user가 들어있는지 확인해줄 함수
func FindUnfollwer(user models.GithubUserInfo, m map[string]int) int {
	val := m[user.Login]
	if val != 1 {
		return 1
	}
	return 0
}

func HitURL(userName string, follow string, i int, c chan models.GithubUserInfo) {
	pageURL := baseurl + userName + "/" + follow + "?per_page=100&page=" + strconv.Itoa(i)
	req, err := http.NewRequest("GET", pageURL, nil)
	utils.CheckErr(err)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	res, err := client.Do(req)
	utils.CheckErr(err)
	body, err := ioutil.ReadAll(res.Body)
	var following User
	err = json.Unmarshal(body, &following)
	utils.CheckErr(err)
	for _, user := range following {
		c <- user
	}
}
