package repository

import (
	"bytes"
	"context"
	"dota-discord-bot/src/internal/kungdota"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type KungdotaRepository interface {
	GetByNames(ctx context.Context, ids []string) (kungdota.Players2, error)
	GetByDiscordID(ctx context.Context, ids []string) ([]kungdota.Players2, error)
	PostMatch(m kungdota.Match) error
	GetAllPlayers() (kungdota.Players2, error)
	SignUp(ctx context.Context, username string, i int) (map[string][]string, error)
	Update(ctx context.Context, username string) (map[string][]string, error)
}

type kungdotaRepository struct {
	config     config
	logger     *zap.Logger
	httpClient http.Client
}

type config struct {
	leagueID string
}

func NewConfig(leagueID string) config {
	return config{
		leagueID: leagueID,
	}
}

func NewKungdotaRepository(logger *zap.Logger, httpClient http.Client, config config) KungdotaRepository {
	return &kungdotaRepository{
		config:     config,
		logger:     logger,
		httpClient: httpClient,
	}
}

func (rx kungdotaRepository) PostMatch(m kungdota.Match) error {
	asd, _ := json.Marshal(m)
	fmt.Println(m)
	buff := bytes.NewBuffer(asd)
	fmt.Println(buff)
	tmp, err := rx.httpClient.Post("https://api.bollsvenskan.jacobadlers.com/match/", "application/json", buff)
	fmt.Println(tmp)

	return err
}

func (rx kungdotaRepository) GetAllPlayers() (kungdota.Players2, error) {
	resp, err := rx.httpClient.Get("https://api.bollsvenskan.jacobadlers.com/player")
	if err != nil {
		rx.logger.Error("KungdotaRepository.GetAllPlayers failed")
		return kungdota.Players2{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.GetAllPlayers failed")
		return kungdota.Players2{}, err
	}

	var players = kungdota.Players2{}
	if err := json.Unmarshal(body, &players); err != nil {
		rx.logger.Error("KungdotaRepository.GetAllPlayers failed")
		return kungdota.Players2{}, err
	}

	return players, nil
}

func (rx kungdotaRepository) GetByNames(ctx context.Context, ids []string) (kungdota.Players2, error) {
	players, err := rx.GetAllPlayers()
	if err != nil {
		rx.logger.Error("KungdotaRepository.GetByNames failed")
		return kungdota.Players2{}, err
	}

	pList := make([]kungdota.Players, 0)
	for _, id := range ids {
		for _, p := range players.Players {
			if p.Username == id {
				pList = append(pList, p)
				goto found
			}
		}
		rx.logger.Error("KungdotaRepository.GetByNames player not found")
		return kungdota.Players2{}, errors.New("player not found")
	found:
	}

	return kungdota.Players2{
		Players: pList,
	}, nil
}

func (rx kungdotaRepository) GetByDiscordID(ctx context.Context, ids []string) ([]kungdota.Players2, error) {
	return nil, nil
}

type signupURL struct {
	SignupDocumentUrl string `json:"signupDocumentUrl"`
	CurrentPollUrl    string `json:"currentPollUrl"`
}

func (rx kungdotaRepository) getURL() (signupURL, error) {
	resp, err := rx.httpClient.Get("https://api.bollsvenskan.jacobadlers.com/dota/signup")
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return signupURL{}, err
	}

	// TODO:
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return signupURL{}, err
	}

	var urls = signupURL{}
	if err := json.Unmarshal(body, &urls); err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return signupURL{}, err
	}

	return urls, err
}

type Csrf struct {
	Token string `json:"token"`
}

func (rx kungdotaRepository) getCSRFToken() (Csrf, error) {
	resp, err := rx.httpClient.Get("https://nextcloud.jacobadlers.com/index.php/csrftoken")
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Csrf{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Csrf{}, err
	}

	var csrf = Csrf{}
	if err := json.Unmarshal(body, &csrf); err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Csrf{}, err
	}

	return csrf, nil
}

type Poll struct {
	Options []Options `json:"options"`
}
type Options struct {
	ID        int       `json:"id"`
	PollID    int       `json:"pollId"`
	Text      time.Time `json:"text"`
	Timestamp int       `json:"timestamp"`
	Order     int       `json:"order"`
	Confirmed int       `json:"confirmed"`
	Duration  int       `json:"duration"`
	Computed  Computed  `json:"computed"`
	Owner     Owner     `json:"owner"`
}
type Computed struct {
	Rank       int  `json:"rank"`
	No         int  `json:"no"`
	Yes        int  `json:"yes"`
	Maybe      int  `json:"maybe"`
	RealNo     int  `json:"realNo"`
	Votes      int  `json:"votes"`
	IsBookedUp bool `json:"isBookedUp"`
}

type Owner struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	IsNoUser    bool   `json:"isNoUser"`
}

func (rx kungdotaRepository) getGameIDs(CSRFToken string, url string) (Poll, error) {
	req, err := http.NewRequest("GET", url+"/options", nil)
	if err != nil {
		fmt.Println(err)
		return Poll{}, err
	}

	req.Header.Add("requesttoken", CSRFToken)

	res, err := rx.httpClient.Do(req)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Poll{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Poll{}, err
	}

	var poll = Poll{}
	if err := json.Unmarshal(body, &poll); err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Poll{}, err
	}

	return poll, nil
}

type CheckUsername struct {
	Token    string `json:"token"`
	UserName string `json:"userName"`
}

type Res struct {
	Result bool   `json:"result"`
	Name   string `json:"name"`
}

func (rx kungdotaRepository) checkUsername(token string, username string, csrf string) (Res, error) {
	payload := &CheckUsername{
		Token:    token,
		UserName: username,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	req, err := http.NewRequest("POST", "https://nextcloud.jacobadlers.com/index.php/apps/polls/check/username", payloadBuf)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Res{}, err
	}

	req.Header.Add("requesttoken", csrf)
	req.Header.Add("Content-Type", "application/json")

	res, err := rx.httpClient.Do(req)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Res{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Res{}, err
	}

	var ret = Res{}
	if err := json.Unmarshal(body, &ret); err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return Res{}, err
	}

	return ret, nil
}

type RegisterUsernameResponse struct {
	Share Share `json:"share"`
}
type Share struct {
	ID              int    `json:"id"`
	Token           string `json:"token"`
	Type            string `json:"type"`
	PollID          int    `json:"pollId"`
	UserID          string `json:"userId"`
	EmailAddress    string `json:"emailAddress"`
	InvitationSent  string `json:"invitationSent"`
	ReminderSent    string `json:"reminderSent"`
	DisplayName     string `json:"displayName"`
	IsNoUser        bool   `json:"isNoUser"`
	URL             string `json:"URL"`
	ShowLogin       bool   `json:"showLogin"`
	PublicPollEmail string `json:"publicPollEmail"`
}

func (rx kungdotaRepository) registerUsername(token string, username string, url string, csrf string) (RegisterUsernameResponse, error) {
	type Reg struct {
		Username     string `json:"userName"`
		EmailAddress string `json:"emailAddress"`
	}

	payload := &Reg{
		Username:     username,
		EmailAddress: "",
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	req, err := http.NewRequest("POST", url+"/register", payloadBuf)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return RegisterUsernameResponse{}, err
	}

	req.Header.Add("requesttoken", csrf)
	req.Header.Add("Content-Type", "application/json")

	res, err := rx.httpClient.Do(req)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return RegisterUsernameResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return RegisterUsernameResponse{}, err
	}

	var ret = RegisterUsernameResponse{}
	if err := json.Unmarshal(body, &ret); err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return RegisterUsernameResponse{}, err
	}

	return ret, nil
}

type Sign struct {
	OptionId int    `json:"optionId"`
	SetTo    string `json:"setTo"`
}

type SignResponse struct {
	Vote    Vote   `json:"vote"`
	Message string `json:"message"`
}
type Vote struct {
	ID         int       `json:"id"`
	PollID     int       `json:"pollId"`
	OptionText time.Time `json:"optionText"`
	Answer     string    `json:"answer"`
	User       User      `json:"user"`
}

func (rx kungdotaRepository) Sign(ID int, token string, csrf string) (SignResponse, error) {
	payload := &Sign{
		OptionId: ID,
		SetTo:    "yes",
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	req, err := http.NewRequest("PUT", "https://nextcloud.jacobadlers.com/index.php/apps/polls/s/"+token+"/vote", payloadBuf)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return SignResponse{}, err
	}

	req.Header.Add("requesttoken", csrf)
	req.Header.Add("Content-Type", "application/json")

	res, err := rx.httpClient.Do(req)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return SignResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return SignResponse{}, err
	}

	var ret = SignResponse{}
	if err := json.Unmarshal(body, &ret); err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return SignResponse{}, err
	}

	return ret, nil
}

type List struct {
	Votes []Votes `json:"votes"`
}
type User struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	IsNoUser    bool   `json:"isNoUser"`
}
type Votes struct {
	ID         int       `json:"id"`
	PollID     int       `json:"pollId"`
	OptionText time.Time `json:"optionText"`
	Answer     string    `json:"answer"`
	User       User      `json:"user"`
}

func (rx kungdotaRepository) updateList(token string, csrf string) (List, error) {
	req, err := http.NewRequest("GET", "https://nextcloud.jacobadlers.com/index.php/apps/polls/s/"+token+"/votes", nil)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return List{}, err
	}

	req.Header.Add("requesttoken", csrf)
	req.Header.Add("Content-Type", "application/json")

	res, err := rx.httpClient.Do(req)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return List{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return List{}, err
	}

	var ret = List{}
	if err := json.Unmarshal(body, &ret); err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return List{}, err
	}

	return ret, nil
}

func (rx List) format() map[string][]string {
	ret := map[string][]string{}

	for _, o := range rx.Votes {
		if o.Answer == "yes" {
			ret[o.OptionText.Format("1504")] = append(ret[o.OptionText.Format("1504")], o.User.DisplayName)
		}
	}

	return ret
}

func (rx kungdotaRepository) SignUp(ctx context.Context, username string, i int) (map[string][]string, error) {
	urls, err := rx.getURL()
	if err != nil {
		return nil, err
	}

	csrf, err := rx.getCSRFToken()
	if err != nil {
		return nil, err
	}

	games, err := rx.getGameIDs(csrf.Token, urls.CurrentPollUrl)
	if err != nil {
		return nil, err
	}

	// TODO:
	// gameTimes := make([]time.Time, 0)
	// for _, o := range games.Options {
	// 	gameTimes = append(gameTimes, time.Unix(int64(o.Timestamp), 0))
	// }

	currentPollToken := strings.ReplaceAll(urls.CurrentPollUrl, "https://nextcloud.jacobadlers.com/index.php/apps/polls/s/", "")

	ok, err := rx.checkUsername(currentPollToken, username, csrf.Token)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp checkUsername failed")
		return nil, err
	}
	if !ok.Result {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return nil, errors.New("username already taken")
	}

	regUserRes, err := rx.registerUsername(currentPollToken, username, urls.CurrentPollUrl, csrf.Token)
	if err != nil {
		return nil, err
	}

	sign, err := rx.Sign(games.Options[i].ID, regUserRes.Share.Token, csrf.Token)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return nil, err
	}
	if sign.Message != "" {
		return nil, errors.New(sign.Message)
	}

	list, err := rx.updateList(currentPollToken, csrf.Token)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return nil, err
	}
	return list.format(), nil
}

func (rx kungdotaRepository) Update(ctx context.Context, username string) (map[string][]string, error) {
	urls, err := rx.getURL()
	if err != nil {
		return nil, err
	}

	csrf, err := rx.getCSRFToken()
	if err != nil {
		return nil, err
	}

	currentPollToken := strings.ReplaceAll(urls.CurrentPollUrl, "https://nextcloud.jacobadlers.com/index.php/apps/polls/s/", "")

	list, err := rx.updateList(currentPollToken, csrf.Token)
	if err != nil {
		rx.logger.Error("KungdotaRepository.SignUp failed")
		return nil, err
	}
	return list.format(), nil
}
