package hooks

import (
	"github.com/devingen/api-core/model"
	"net/http"
	"os"
	"time"
)

const (
	envWebHookAddress = "DATA_API_WEB_HOOK_ADDRESS"
)

var hookURL string

func init() {
	hookURL, _ = os.LookupEnv(envWebHookAddress)
}

func CheckEligibility(authorizationHeader string) error {

	if hookURL == "" {
		return nil
	}
	request, err := http.NewRequest(http.MethodPost, hookURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", authorizationHeader)

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	isEligible := resp.StatusCode == http.StatusOK

	if !isEligible {
		return model.NewErrorWithCode(http.StatusUnauthorized, 0, "web-hook-declined")
	}
	return nil
}
