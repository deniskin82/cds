package test

import (
	"encoding/base64"
	"io"
	"net/http"
	"testing"

	"github.com/ovh/cds/engine/test/config"
	"github.com/stretchr/testify/require"
)

// TestKey is a test key encoded in base64, do not use it in real life
var TestKey []byte

func init() {
	TestKey, _ = base64.StdEncoding.DecodeString(`LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUNGd0FBQUFkemMyZ3RjbgpOaEFBQUFBd0VBQVFBQUFnRUFwNDN4L2U3S3B6TjJLMDc2Y1FhU3ppL3B0YTQyMVYvMnNPRGpuUE52dWpkamlQaEtYWW1kClFmNmtIZHFUeTVKdW5KVkJ5cTFFM0t6dWd0cCs1VUlGZ21tYVNWU092R1B0MkdZaERiemVsSUlsWTFhVlRvQ3BnL2hTb3gKRFhucXBlMkFRUjRIU1hHdEJTQUJVbXJYb1ZxNGRtZmFtYXloVzc5NmpTVEdMY0lMWFBJNXFUcFNlMjd3cGNGOGM5MXFGVgpUNFo5b0FMbzB2ZmlYQlF6c0tlMVNBbjZlOWtGUEtJZ01wak91Y3kyTkVuVDVLTTlJVFNSOTRCM0ptcDV5enBOMXpHdmRMCkxTQ2ptUGlxQnpCSUlncCt1T01WMDBZT2JiUHJOVTM5NHJhYzJWR2tOb2JYb3kwL05YUUxQY0hGeUxSMXoyVGxsWUNabWUKUVZKV1laelVtYm5CaWN4SXlDYVQ1bExGdmorTjZkbWd5cE90aTJXSEY1MHJlOG0xa0ZDTll4UlkrWjAzTDliNzFiS1UzcgpzOVM2RThHMDVZWVFQOER6MWdzUWRvSVRwT0w4cnROaDd5SWdCaG8weVozQ2RnQ09Pem1rM0hOaU5NQXdNUS9NWHZXMytLCjNhTFMxMG1tb3dJT0ZXdWtKT2cxWkJBRWVndnpNZ2xxcGZ1eElWNWljT3oxOUVST0Q1bkZjeGl5TkV0dTAza3hTVzhnRjUKYlBNTnFydkNFNytqSFVxSzlhRmpxeHFtOWZvSjJRZnlZcG5qQmZpd1JkWG9yZzBwSzcvRnNYeFhHMkZ6ZHF4RmNtK2FEbwpGbnlxZzExWTFhdlBYKzZ2WXk5bE4rRWRNQ1NNVUt3U3NSMkdaMkg1M1FDczFSQWhsaG10MDFtS2k2Y0Q1VTVLRVJGTkxWClVBQUFkZ1pyWlVLV2EyVkNrQUFBQUhjM05vTFhKellRQUFBZ0VBcDQzeC9lN0twek4ySzA3NmNRYVN6aS9wdGE0MjFWLzIKc09Eam5QTnZ1amRqaVBoS1hZbWRRZjZrSGRxVHk1SnVuSlZCeXExRTNLenVndHArNVVJRmdtbWFTVlNPdkdQdDJHWWhEYgp6ZWxJSWxZMWFWVG9DcGcvaFNveERYbnFwZTJBUVI0SFNYR3RCU0FCVW1yWG9WcTRkbWZhbWF5aFc3OTZqU1RHTGNJTFhQCkk1cVRwU2UyN3dwY0Y4YzkxcUZWVDRaOW9BTG8wdmZpWEJRenNLZTFTQW42ZTlrRlBLSWdNcGpPdWN5Mk5FblQ1S005SVQKU1I5NEIzSm1wNXl6cE4xekd2ZExMU0NqbVBpcUJ6QklJZ3ArdU9NVjAwWU9iYlByTlUzOTRyYWMyVkdrTm9iWG95MC9OWApRTFBjSEZ5TFIxejJUbGxZQ1ptZVFWSldZWnpVbWJuQmljeEl5Q2FUNWxMRnZqK042ZG1neXBPdGkyV0hGNTByZThtMWtGCkNOWXhSWStaMDNMOWI3MWJLVTNyczlTNkU4RzA1WVlRUDhEejFnc1Fkb0lUcE9MOHJ0Tmg3eUlnQmhvMHlaM0NkZ0NPT3oKbWszSE5pTk1Bd01RL01YdlczK0szYUxTMTBtbW93SU9GV3VrSk9nMVpCQUVlZ3Z6TWdscXBmdXhJVjVpY096MTlFUk9ENQpuRmN4aXlORXR1MDNreFNXOGdGNWJQTU5xcnZDRTcrakhVcUs5YUZqcXhxbTlmb0oyUWZ5WXBuakJmaXdSZFhvcmcwcEs3Ci9Gc1h4WEcyRnpkcXhGY20rYURvRm55cWcxMVkxYXZQWCs2dll5OWxOK0VkTUNTTVVLd1NzUjJHWjJINTNRQ3MxUkFobGgKbXQwMW1LaTZjRDVVNUtFUkZOTFZVQUFBQURBUUFCQUFBQ0FRQ1FNU0NTeGdBU09jQTA3d2VwWXQzTm9RQUFRTWVoZ3E4YQpjcjZPWUJURGJVMDBIM0JuNUxpM2hYc1kwZlNrbVFTbHJmRHJpWWNjWFpuNGRDNEYvM1ljVCtMZHZtNERnLys0WGRPT0xmCjVpVVVuNW5oWnBjMkh1VnpKT2NIME9aMUd0bG5zSDdXM29QbVNDKzdESVU2cjRiVkp2VEJrUVZmbm4zSm4xOEpHOWVKaWsKN0M2cFQyOG5jWVBsVnFwSjNaYzhFK0ppWkg2V3A0cGVjV2cyVzIwdmJKN3FHODVjNnF6SXZpWVJVVEZ2K0NUb3V1NHRlRAo4eGZwV0xNdEJUYTM1M2RhT255d2Zra3JxTHN4Nm9QNC80MGtjUkJrUEFMSXQ2L3Z0SW1McEZtQXo3aUEwRFFja2lDMlVJCklvQ0d5OEYwalhUTjRpZFlRNklrVnNaTnhKaFR1YkdCeUg2V0JVdjV3VlNxd2FQUk5yWTRtRE1ONUI2N04vWEZIdjVhejcKQjUwWUVwYUhKTmxTdlphc2JWSmQ5QU0vU0x0cERuNDArZ0gvUzdXVDVvTTRpQWtUcWxzdjZRRHQzTURJOVZ5NTluaGFsMgpNQXdZc3NzbnZldWNuYVR5OUFOQWdHNjMzb3E1N1BReGFGdEFVUkFxQ3ptNVhhbjY2UTVISkhKZHJQcWdXM3RrRzNxall6CmFoMnFGdW05bHFhMTJtWVRNdTZkNlRiUmtFQ1lkTk02QTRwdnIxeVhWQkY0M09ncjU0T1FXaVNKV0Jka3JObHNKcGl5RUwKcjlRSG9nNkpIRmc0ckJWT0xjZDkvVWV0R3NDSnpTZ0JNNWxwZFV0RmdlVGJVREUzSDhHbjdqZlZRd1VYcW8xVENjdnh6bAp5ZWdRaXFTSThtaHJxMzJUNllkUUFBQVFCNlY2a2gxSjA2ZVlveFBYMjZCNlIrVC9XWUM3Tkp3eXdLSzJjTE5FbjBDWE9pClFqWnFKOHdocWo0TVl2Q2xmTGJSNUtNd2dib3J5NU1yQXJpQ2NzVVRTRy9XSWlFTlJEM005VUo3ZnN0bGY5bWllQlNWKzUKVEovMkZ3VTFabU82ci9nbC9BNFRNMGJUVXR1aFJZc0x0SFZVS3pFaW5LVVpUQi83S3l1SWljNDhIL1ZNcFNUdlpuUzVSagpRVzQvN0dUQ2tFSERMUEJnNmJZSGVPM2YrT1hkTTJnSnlaL3JrTmwvVm1RSGIrYS84Z2FUL21uTEtGNDhBVlkxQmUxMFNTCkoxcHM5eEQ5bm1IMWJZMGFlUUo3dEZkTE1IZGhoV0Y3Q05lT21xSE5ZMVdMaHczeC90cEZmcVJjbTYyTUp6ODhhaFNickwKekxRbWVONVR0U2hYeHJlK0FBQUJBUURiVWxyWG9XeXRHOHFGWUVqeGMrQmZFM3ltYTJwVlpCWjg2TEtnRjk5Rmo0dkdkRApMNWd3dDBGMTNLblpndUZka0UycHJGczlrbmxGQlFMbHcyZWRiVnVhUERrcURncGdVOFlEaEdxSXNzUlZpR0hCNmJqZ1hpClFaeW1LM1B3K005eERKVmpEYnIyVUt2Wk5ZWnVUWkFvcGVTVkVyaW84bFJsVUw2d1dOMm9RMTVGcWJYVysvSzVXaStXWHgKb3lEZTlyVGFZKzU1UndFMTNHT2JlNG5ZanJmeVl5OHhaZGthU3FWMWJVU2NlVkNtaXc3WnJ2eUs3TDVIRnh3ME11eUg2RApsMlZCMElZL1RZU21OYlhyZVpnR21uUkFJbmk5VzFyS0J2K3JXWmllVzV3RzdxdlNlSm9NamU4aGNmeXZRWHFTakxXRlR4CjNWY2xybjVyMUo4bFQzQUFBQkFRRERrMUpNQ0xiQVpIYVVBYmh2MGVxSFlLRWFpWTBnKzZZakNmSnNoMi96eEs2NnZiYUQKak5NTmdxQXNqem4zUU9Kb1hBbmRCRmxKOTZqKzNqQlpaa3dwMVFyY2N3eUtrUUI1NGdXRjV2NHFXSjZjNkZSZUFpYzR4bwo1a2l4cTNqVUJNcG85TDNIQ3dGYlpxVzZqTkI3NDlHSXo2VEhCdEFZYWdaWHpMZ2haTDJLcGhqM2c0SXkzaW9TTkk3NFc1CllzQXFLSWhxNGFFTitJTERIdjhuaTRaWWM3OXpFNGJXRCtTRmp1aC9XazJweDJUM2RIaHRLaGg4eXRhamh0UlNhcVowY2MKRlIrSy93OE1PVDVzelpIbllzaVV4R1BnKzkvZVlZYlZ6azZGazJHL2xrN0RXZ1BweitPVi9QS29QMGs5RDhUZEZzbE5EWApvODBOYmdxYXNWa1RBQUFBS0daeVlXNWpiMmx6TG5OaGJXbHVLMmRwZEdoMVltWmhhMlZyWlhsQVpYaGhiWEJzWlM1amIyCjBCQWc9PQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0K`)
}

// LoadTestingConf loads test configuration tests.cfg.json
func LoadTestingConf(t require.TestingT, serviceType string) map[string]string {
	return config.LoadTestingConf(t, serviceType)
}

//GetTestName returns the name the test
func GetTestName(t *testing.T) string {
	return t.Name()
}

//FakeHTTPClient implements sdk.HTTPClient and returns always the same response
type FakeHTTPClient struct {
	T        *testing.T
	Response *http.Response
	Error    error
}

//Do implements sdk.HTTPClient and returns always the same response
func (f *FakeHTTPClient) Do(r *http.Request) (*http.Response, error) {
	b, err := io.ReadAll(r.Body)
	if err == nil {
		r.Body.Close()
	}

	f.T.Logf("FakeHTTPClient> Do> %s %s: Payload %s", r.Method, r.URL.String(), string(b))
	return f.Response, f.Error
}
