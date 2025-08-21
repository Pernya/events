package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ic2hrmk/promtail"
)

type VK_Segments struct {
	Count int `json:"count"`
	Items []struct {
		Base            int    `json:"base"`
		Created         string `json:"created"`
		EntriesCount    int    `json:"entries_count"`
		HasHistory      bool   `json:"has_history"`
		Id              int    `json:"id"`
		IdsCount        int    `json:"ids_count"`
		MatchedIdsCount int    `json:"matched_ids_count"`
		Name            string `json:"name"`
		Status          string `json:"status"`
		Type            string `json:"type"`
	} `json:"items"`
	Limit  interface{} `json:"limit"`
	Offset int         `json:"offset"`
}
type YaSegment struct {
	ID                         int    `json:"id"`
	Name                       string `json:"name"`
	Type                       string `json:"type"`
	Status                     string `json:"status"`
	CreateTime                 string `json:"create_time"`
	Owner                      string `json:"owner"`
	HasGuests                  bool   `json:"has_guests"`
	GuestQuantity              int    `json:"guest_quantity"`
	CanCreateDependent         bool   `json:"can_create_dependent"`
	HasDerivatives             bool   `json:"has_derivatives"`
	CookiesMatchedQuantity     int    `json:"cookies_matched_quantity"`
	Hashed                     bool   `json:"hashed,omitempty"`
	UsedHashingAlg             string `json:"used_hashing_alg,omitempty"`
	ContentType                string `json:"content_type,omitempty"`
	ItemQuantity               int    `json:"item_quantity,omitempty"`
	ValidUniqueQuantity        int    `json:"valid_unique_quantity,omitempty"`
	ValidUniquePercentage      string `json:"valid_unique_percentage,omitempty"`
	MatchedQuantity            int    `json:"matched_quantity,omitempty"`
	MatchedPercentage          string `json:"matched_percentage,omitempty"`
	CounterID                  int    `json:"counter_id,omitempty"`
	UploadingLastModifyTime    string `json:"uploading_last_modify_time,omitempty"`
	Guest                      bool   `json:"guest"`
	LookalikeLink              int    `json:"lookalike_link,omitempty"`
	LookalikeValue             int    `json:"lookalike_value,omitempty"`
	MaintainDeviceDistribution bool   `json:"maintain_device_distribution,omitempty"`
	MaintainGeoDistribution    bool   `json:"maintain_geo_distribution,omitempty"`
	Derivatives                []int  `json:"derivatives,omitempty"`
	PixelID                    int    `json:"pixel_id,omitempty"`
	PeriodLength               int    `json:"period_length,omitempty"`
}
type YaSegments struct {
	Segments []struct {
		ID                         int    `json:"id"`
		Name                       string `json:"name"`
		Type                       string `json:"type"`
		Status                     string `json:"status"`
		CreateTime                 string `json:"create_time"`
		Owner                      string `json:"owner"`
		HasGuests                  bool   `json:"has_guests"`
		GuestQuantity              int    `json:"guest_quantity"`
		CanCreateDependent         bool   `json:"can_create_dependent"`
		HasDerivatives             bool   `json:"has_derivatives"`
		CookiesMatchedQuantity     int    `json:"cookies_matched_quantity"`
		Hashed                     bool   `json:"hashed,omitempty"`
		UsedHashingAlg             string `json:"used_hashing_alg,omitempty"`
		ContentType                string `json:"content_type,omitempty"`
		ItemQuantity               int    `json:"item_quantity,omitempty"`
		ValidUniqueQuantity        int    `json:"valid_unique_quantity,omitempty"`
		ValidUniquePercentage      string `json:"valid_unique_percentage,omitempty"`
		MatchedQuantity            int    `json:"matched_quantity,omitempty"`
		MatchedPercentage          string `json:"matched_percentage,omitempty"`
		CounterID                  int    `json:"counter_id,omitempty"`
		UploadingLastModifyTime    string `json:"uploading_last_modify_time,omitempty"`
		Guest                      bool   `json:"guest"`
		LookalikeLink              int    `json:"lookalike_link,omitempty"`
		LookalikeValue             int    `json:"lookalike_value,omitempty"`
		MaintainDeviceDistribution bool   `json:"maintain_device_distribution,omitempty"`
		MaintainGeoDistribution    bool   `json:"maintain_geo_distribution,omitempty"`
		Derivatives                []int  `json:"derivatives,omitempty"`
		PixelID                    int    `json:"pixel_id,omitempty"`
		PeriodLength               int    `json:"period_length,omitempty"`
	} `json:"segments"`
}

const vk_config_filename = "vk_config.json"

type TokenResponseVK struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (c *TokenResponseVK) getUserList()

type YaSegmentRequestSingle struct {
	Segment struct {
		ID          int    `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Hashed      bool   `json:"hashed,omitempty"`
		HashingAlg  string `json:"hashing_alg,omitempty"`
		ContentType string `json:"content_type,omitempty"`

		LookalikeLink              int  `json:"lookalike_link,omitempty"`
		LookalikeValue             int  `json:"lookalike_value,omitempty"`
		MaintainDeviceDistribution bool `json:"maintain_device_distribution,omitempty"`
		MaintainGeoDistribution    bool `json:"maintain_geo_distribution,omitempty"`
	} `json:"segment"`
}
type YaSegmentResponseSingle struct {
	Segment struct {
		ID                 int    `json:"id"`
		Name               string `json:"name,omitempty"`
		Type               string `json:"type"`
		Status             string `json:"status"`
		HasGuests          bool   `json:"has_guests"`
		GuestQuantity      int    `json:"guest_quantity"`
		CanCreateDependent bool   `json:"can_create_dependent"`
		HasDerivatives     bool   `json:"has_derivatives"`
		Hashed             bool   `json:"hashed"`
		UsedHashingAlg     string `json:"used_hashing_alg"`
		ItemQuantity       int    `json:"item_quantity"`
		Guest              bool   `json:"guest"`
		ContentType        string `json:"content_type,omitempty"`
	} `json:"segment"`
}
type YandexApp struct {
	Token    string `json:"token"`
	ClientId string `json:"client_id"`
	FullUrl  string `json:"full_url"`
}
type ConfigFile struct {
	Token    string `json:"token"`
	ClientID string `json:"client_id"`
	FullURL  string `json:"full_url"`
}

func (c *ConfigFile) CreateSegment(method string, url_full string, data interface{}) (string, string)
func (c *ConfigFile) GetSegments(method string, url_full string) (string, string)
func (c *ConfigFile) UploadFile(url string, paramName string, filePath string) (string, string)

type SegmentFileAccelera struct {
	Filename    string    `json:"filename"`
	FileContent []string  `json:"file_content"`
	SegmentName string    `json:"segment_name"`
	Timestamp   time.Time `json:"timestamp"`
}

func (s SegmentFileAccelera) FileNameSetter(segment_name string) SegmentFileAccelera
func (s SegmentFileAccelera) PopulateSegmentMap(sma map[string]SegmentMapAccelera, data_full []string, segment_name string) map[string]SegmentMapAccelera
func (s SegmentFileAccelera) ReadSegmentFromFile(segment_full_map map[string]SegmentMapAccelera) ([]string, error)

type SegmentMapAccelera struct {
	Columns     map[string]string
	FieldValues map[string]bool
	SegmentName string
	Timestamp   time.Time
}

func main() {
	r := gin.Default()
	segment_full_map := make(map[string]SegmentMapAccelera)
	r.Use(cors.Default())
	f_yandex, _ := os.Open("ya_auth.json")

	f_content, _ := io.ReadAll(f_yandex)
	var configFile ConfigFile
	log.SetOutput(gin.DefaultWriter)
	fmt.Println(string(f_content))

	err_config := json.Unmarshal(f_content, &configFile)
	if err_config != nil {

		fmt.Errorf(err_config.Error())

	}
	r.Use(gin.Recovery())
	r.POST("/post", func(c *gin.Context) {

		ya_url := c.PostForm("url")

		//test_url := `https://lb.4lapy.accelera.ai/auth/yandex#access_token=y0__wgBEMzY_JUGGLCANCCm04XrERuS0Yo-SNEcWQuP9GU3WxEqoMmb&token_type=bearer&expires_in=31536000&cid=00mm24d3k9yzb5xfv18dv239tc`
		//
		//fmt.Println("ya_url", strings.Replace(ya_url, "#", "?", 1))
		myUrl, _ := url.Parse(strings.Replace(ya_url, "#", "?", 1))
		params, _ := url.ParseQuery(myUrl.RawQuery)
		//
		fmt.Println("Токен:", params["access_token"])

		f, _ := os.Create("ya_auth.json")

		ya_app := YandexApp{
			Token:    params["access_token"][0],
			ClientId: "",
			FullUrl:  ya_url,
		}

		ya_app_str, _ := json.Marshal(ya_app)

		_, err_write := f.WriteString(string(ya_app_str))
		if err_write != nil {
			fmt.Errorf("ошибка записи в файл %s", err_write.Error())
		}

		f.Close()
		//fmt.Println(ya_url)

		c.HTML(http.StatusOK, "thank_you_for_reg.tmpl", gin.H{})
	})

	r.POST("/update_segment/yandex", func(c *gin.Context) {

		segment_name := c.DefaultQuery("segment_name", "")

		// TIP ЭТО ID АУДИТОРИИ ЯНДЕКСА!!!!

		audience_id := c.DefaultQuery("audience_id", "")
		RefreshSegmentAcceleraHandler(c, segment_full_map)

		os.Open("segments/" + segment_name + ".csv")
		var segment SegmentFileAccelera

		segment.SegmentName = segment_name
		status, response := configFile.UploadFile(fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%s/modify_data?modification_type=replace", audience_id), "file", "segments/"+segment_name+".csv")

		fmt.Println(fmt.Sprintf("Код от сервера %s", status))
		fmt.Println(fmt.Sprintf("Ответ от сервера %s", response))

		if status != "200" {

			RecordLoki(map[string]string{

				"events":       "yandex_audience",
				"event_status": "error",
				"event_name":   "create_error",
				"segment_name": segment_name,
			}, "error", "error happened ")

		}

		var YaResponse YaSegmentResponseSingle
		var Ya_request YaSegmentRequestSingle
		err_json := json.Unmarshal([]byte(response), &YaResponse)
		if err_json != nil {

			fmt.Println(err_json)

		}

		ya_resp_json, _ := json.Marshal(YaResponse)

		Ya_request.Segment.ContentType = "crm"
		Ya_request.Segment.Name = YaResponse.Segment.Name
		Ya_request.Segment.Hashed = YaResponse.Segment.Hashed
		Ya_request.Segment.HashingAlg = YaResponse.Segment.UsedHashingAlg

		fmt.Println("Получен json для отправки " + string(ya_resp_json))

		seg_dat, _ := json.Marshal(Ya_request)

		status2, response2 := configFile.CreateSegment("post", fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%s/confirm", strconv.Itoa(YaResponse.Segment.ID)), string(seg_dat))

		fmt.Println(fmt.Sprintf("Код от сервера %s", status2))
		fmt.Println(fmt.Sprintf("Ответ от сервера %s", response2))

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", "segments/"+segment_name+".csv"))
	})
	r.POST("/update_segment/vk", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
	r.POST("/lookalike_segment/yandex", func(c *gin.Context) {

		segment_name := c.DefaultQuery("segment_name", "")
		title := c.DefaultQuery("title", "")
		// TIP ЭТО ID АУДИТОРИИ ЯНДЕКСА!!!!

		audience_id := c.DefaultQuery("audience_id", "")
		RefreshSegmentAcceleraHandler(c, segment_full_map)

		os.Open("segments/" + segment_name + ".csv")
		//var segment SegmentFileAccelera

		//		var YaResponse YaSegmentResponseSingle
		var Ya_request YaSegmentRequestSingle
		//err_json := json.Unmarshal([]byte(response), &YaResponse)
		//if err_json != nil {
		//
		//	fmt.Println(err_json)
		//
		//}
		//
		//ya_resp_json, _ := json.Marshal(YaResponse)

		a_i, err_id := strconv.Atoi(audience_id)

		if err_id != nil {

			c.String(http.StatusBadRequest, "audience_id is malformed or corrupt")

		}

		Ya_request.Segment.ContentType = "crm"
		Ya_request.Segment.Name = title
		Ya_request.Segment.LookalikeLink = a_i
		Ya_request.Segment.LookalikeValue = 5
		Ya_request.Segment.MaintainGeoDistribution = false
		Ya_request.Segment.MaintainDeviceDistribution = false

		//Ya_request.Segment.Hashed = YaResponse.Segment.Hashed
		//Ya_request.Segment.HashingAlg = YaResponse.Segment.UsedHashingAlg

		ya_req, _ := json.Marshal(Ya_request)

		status, response := configFile.CreateSegment("post", "https://api-audience.yandex.ru/v1/management/segments/create_lookalike", string(ya_req))

		fmt.Println(fmt.Sprintf("Код от сервера %s", status))
		fmt.Println(fmt.Sprintf("Ответ от сервера %s", response))
		//fmt.Println("Получен json для отправки " + string(ya_resp_json))
		var YaResponse YaSegmentResponseSingle

		err_json := json.Unmarshal([]byte(response), &YaResponse)
		if err_json != nil {

			fmt.Println(err_json)

		}
		seg_dat, _ := json.Marshal(Ya_request)

		status2, response2 := configFile.CreateSegment("post", fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%s/confirm", strconv.Itoa(YaResponse.Segment.ID)), string(seg_dat))

		fmt.Println(fmt.Sprintf("Код от сервера %s", status2))
		fmt.Println(fmt.Sprintf("Ответ от сервера %s", response2))

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", title))
	})
	r.POST("/lookalike_segment/vk", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
	r.POST("/create_segment/yandex", func(c *gin.Context) {
		//file, _ := c.FormFile("file")
		//log.Println(file.Filename)
		title := c.DefaultQuery("title", "")
		segment_name := c.DefaultQuery("segment_name", "")

		// Upload the file to specific dst.
		//err_file_save := c.SaveUploadedFile(file, file.Filename)
		//
		//READ:
		//	result, err := segment.ReadSegmentFromFile(segment_full_map)
		//
		//	if err != nil {

		RefreshSegmentAcceleraHandler(c, segment_full_map)

		//os.Open("segments/" + segment_name + ".csv")
		var segment SegmentFileAccelera

		//segment.SegmentName = segment_name
		segment.SegmentName = title
		//
		//	/* и пусть меня проклянут , но я обожаю богоподобный goto */
		//	goto READ
		//}

		//if err_file_save != nil {
		//
		//	c.String(http.StatusBadRequest, "error during file processing: "+err_file_save.Error())
		//
		//}
		max_retries := 3
		var YaResponse YaSegmentResponseSingle
		var Ya_request YaSegmentRequestSingle
		for i := 0; i <= max_retries; i++ {

			status, response := configFile.UploadFile("https://api-audience.yandex.ru/v1/management/segments/upload_csv_file", "file", "segments/"+segment_name+".csv")

			switch status {
			case "200":
				RecordLoki(map[string]string{

					"events":       "yandex_audience",
					"event_status": "success",
					"event_name":   "create_segment",
					"segment_name": segment_name,
					"retry":        strconv.Itoa(i),
				}, "info", response)

				fmt.Println(fmt.Sprintf("Код от сервера %s", status))
				fmt.Println(fmt.Sprintf("Ответ от сервера %s", response))

				err_json := json.Unmarshal([]byte(response), &YaResponse)
				if err_json != nil {

					fmt.Println(err_json)

				}

				ya_resp_json, _ := json.Marshal(YaResponse)

				Ya_request.Segment.ContentType = "crm"
				Ya_request.Segment.Name = title
				Ya_request.Segment.Hashed = YaResponse.Segment.Hashed
				Ya_request.Segment.HashingAlg = YaResponse.Segment.UsedHashingAlg

				fmt.Println("Получен json для отправки " + string(ya_resp_json))

				seg_dat, _ := json.Marshal(Ya_request)

				status2, response2 := configFile.CreateSegment("post", fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%s/confirm", strconv.Itoa(YaResponse.Segment.ID)), string(seg_dat))

				fmt.Println(fmt.Sprintf("Код от сервера %s", status2))
				fmt.Println(fmt.Sprintf("Ответ от сервера %s", response2))

				c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", "segments/"+segment_name+".csv"))

			default:
				RecordLoki(map[string]string{

					"events":       "yandex_audience",
					"event_status": "error",
					"event_name":   "create_error",
					"segment_name": segment_name,
					"retry":        strconv.Itoa(i),
				}, "error", response)

				fmt.Println(fmt.Sprintf("Код от сервера %s", status))
				fmt.Println(fmt.Sprintf("Ответ от сервера %s", response))

				err_json := json.Unmarshal([]byte(response), &YaResponse)
				if err_json != nil {

					fmt.Println(err_json)

				}

				ya_resp_json, _ := json.Marshal(YaResponse)

				Ya_request.Segment.ContentType = "crm"
				Ya_request.Segment.Name = title
				Ya_request.Segment.Hashed = YaResponse.Segment.Hashed
				Ya_request.Segment.HashingAlg = YaResponse.Segment.UsedHashingAlg

				fmt.Println("Получен json для отправки " + string(ya_resp_json))

				seg_dat, _ := json.Marshal(Ya_request)

				status2, response2 := configFile.CreateSegment("post", fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%s/confirm", strconv.Itoa(YaResponse.Segment.ID)), string(seg_dat))

				fmt.Println(fmt.Sprintf("Код от сервера %s", status2))
				fmt.Println(fmt.Sprintf("Ответ от сервера %s", response2))

				//c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", "segments/"+segment_name+".csv"))

			}

		}

	})

	r.GET("/test_record_loki", func(c *gin.Context) {

		var segment_name = "test_segment"

		err_loki := RecordLoki(map[string]string{
			"events":       "yandex_audience",
			"event_status": "error",
			"event_name":   "create_error",
			"segment_name": segment_name},

			"error",

			`{"segment":{"id":45543025,"type":"uploading","status":"uploaded","has_guests":false,"guest_quantity":0,"can_create_dependent":false,"has_derivatives":false,"hashed":false,"used_hashing_alg":"IDENTITY","item_quantity":500,"guest":false}`)

		if err_loki != nil {

			fmt.Println(err_loki)

		}

		c.String(200, "record created")
	})

	r.POST("/create_segment/vk", func(c *gin.Context) {

		title := c.DefaultQuery("title", "")
		//segment_name := c.DefaultQuery("segment_name", "")
		list_type := c.DefaultQuery("list_type", "phones")
		// Upload the file to specific dst.
		//err_file_save := c.SaveUploadedFile(file, file.Filename)
		//
		//READ:
		//	result, err := segment.ReadSegmentFromFile(segment_full_map)
		//
		//	if err != nil {

		RefreshSegmentAcceleraHandler(c, segment_full_map)

		//os.Open("segments/" + segment_name + ".csv")
		VK_CONFIG := TokenResponseVK{}
		err_create := VK_CONFIG.uploadUserListVK("https://ads.vk.com/api/v3/remarketing/users_lists.json", list_type, title, "../vk_config/"+vk_config_filename)

		if err_create != nil {

			RecordLoki(map[string]string{

				"events":       "vk_audience",
				"event_status": "error",
				"event_name":   "create_error",
			}, "error", err_create.Error())

			c.JSON(400, struct {
				ReplyCode int    `json:"reply_code"`
				ReplyText string `json:"reply_text"`
			}{
				ReplyCode: 1,
				ReplyText: "error on segment creation",
			})

		}

		c.JSON(200, gin.H{})
	})
	r.GET("/get_segments/yandex", func(c *gin.Context) {

		statusCode, response := configFile.GetSegments("get", "https://api-audience.yandex.ru/v1/management/segments")

		fmt.Printf("ответ от сервера: %s \n", statusCode)
		//fmt.Printf("ответ от сервера: %s", response)

		var yaSegments YaSegments

		json.Unmarshal([]byte(response), &yaSegments)

		c.JSON(http.StatusOK, yaSegments)
	})
	r.GET("/get_segments/vk", func(c *gin.Context) {

		//response, err_test := makeCurlLikeRequest()
		//fmt.Println(response)
		//fmt.Println(err_test)

		//resp, err_token, vkConf := refreshTokenVK("https://ads.vk.com/api/v2/oauth2/token.json")
		refreshTokenVK("https://ads.vk.com/api/v2/oauth2/token.json")
		resp, err_token := refreshTokenVKCURL()
		fmt.Println("Ответ от VK: ", string(resp))
		var vkConf TokenResponseVK

		json.Unmarshal([]byte(resp), &vkConf)

		//vkConf.uploadUserListVK()

		if err_token != nil {
			fmt.Println(fmt.Errorf("Ошибка получения токена %s", err_token.Error()))
			RecordLoki(map[string]string{

				"events":       "vk_audience",
				"event_status": "error",
				"event_name":   "token_error",
			}, "error", err_token.Error())

			c.JSON(400,

				struct {
					ReplyCode int    `json:"reply_code"`
					ReplyText string `json:"reply_text"`
				}{
					ReplyCode: 1,
					ReplyText: "error on token obtainment",
				})
			return
		}

		VK_CONFIG = vkConf

		c.JSON(200, struct {
			ReplyCode int    `json:"reply_code"`
			ReplyText string `json:"reply_text"`
		}{
			ReplyCode: 0,
			ReplyText: "token obtained successfully",
		})
	})
}
func RecordLoki(identifiers map[string]string, message_type string, message string) error {

	promtailClient, err := promtail.NewJSONv1Client("http://51.250.44.94:3100", identifiers)

	if err != nil {

		fmt.Println(
			err)
		return err
	}
	//defer promtailClient.Close()

	switch message_type {
	case "info":
		promtailClient.Infof(message)
	case "error":
		promtailClient.Errorf(message)
	default:
		promtailClient.Warnf(message)
	}

	return nil

}
func RefreshSegmentAcceleraHandler(c *gin.Context, segment_full_map map[string]SegmentMapAccelera) error {
	//auth_token := c.GetHeader("Authorization")
	segment_name := c.DefaultQuery("segment_name", "")
	//var segment_contents []string

	if segment_name != "" {

		var segmentFile SegmentFileAccelera

		segmentFile.SegmentName = segment_name
		const (
			host     = "rc1d-ne3qaom0jkumgf9i.mdb.yandexcloud.net"
			port     = 6432
			user     = "nwmtod1"
			password = "Gj&d6nf*5R%S!,K"
			dbname   = "analytic_db"
		)
		// connection string
		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s ", host, port, user, password, dbname)

		// open database
		db, err := sql.Open("postgres", psqlconn)
		defer db.Close()

		if err != nil {

			fmt.Println(err)
			return err
		}
		//query := fmt.Sprintf(`SELECT * FROM segments.%s limit 1`, segment_name)
		//rows, err := db.Query(query)
		//if err != nil {
		//
		//	fmt.Println(err)
		//
		//}
		//
		//defer rows.Close()
		////fmt.Println(string(rows))
		//
		////var row string
		//
		//column_names, err_col := rows.Columns()
		//
		//if err_col != nil {
		//
		//	fmt.Println(err_col)
		//}
		//fmt.Println(column_names)

		qeyr_full := fmt.Sprintf("select * from %s", segment_name)

		rows, err_full_query := db.Query(qeyr_full)

		if err_full_query != nil {

			fmt.Println(err_full_query)
			return err_full_query
		}
		defer rows.Close()
		cols, err := rows.Columns()
		rawResult := make([][]byte, len(cols))
		result := make([]string, len(cols))

		dest := make([]interface{}, len(cols)) // A temporary interface{} slice
		for i, _ := range rawResult {
			dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
		}

		var full_data []string
		full_data = append(full_data, strings.Join(cols, ","))
		for rows.Next() {
			err = rows.Scan(dest...)
			if err != nil {
				fmt.Println("Failed to scan row", err)
				return err
			}

			for i, raw := range rawResult {
				if raw == nil {
					result[i] = "null"
				} else {
					result[i] = string(raw)
				}
			}

			//	fmt.Printf("%#v\n", result)

			full_data = append(full_data, strings.Join(result, ","))

		}

		segmentFile.Filename = segment_name + ".csv"
		file, err := os.Create("segments/" + segmentFile.Filename)

		if err != nil {

			return err
		}

		defer file.Close()

		fmt.Println(len(full_data))

		for _, i := range full_data {

			file.WriteString(i + "\n")

		}
		//recieved_data := make(map[string]interface{})

		//for result.NextResultSet() {
		//	var res sql.NullString
		//
		//	result.Scan(&res)
		//
		//	fmt.Println(res)
		//
		//}

	}
	return nil
}
func refreshTokenVK(apiURL string) (string, error, TokenResponseVK) {
	// Подготавливаем данные формы
	//formData := url.Values{
	//	"grant_type":    {"refresh_token"},
	//	"client_id":     {"CyY9AIlj6zRk3nVy"},
	//	"client_secret": {"SmnvMxy3UiWZ89NNrHMtIQsg8C7AMljEsLoJzdJKYs1XrzMOBZmGYfIhKqeoOVDdu44ohHXmEcQGLaFYTw7rePSAfJcU8ToTmNg5yVnfMbgWsYRmcvAvU5sc0dK52oYaKwiT7UHMhM8sZie8OCxguh5ghj9rm0rdaAHZNNJ34bvRSS27Ixa5RYwuhJsG3EBoB5IEmxgT0WnhPLQLKqX5ciokO1JiOUbfJswwypz6AKZhfdZCDRS3wpkNor9"},
	//	"refresh_token": {"4OCe2FGdcneaF6cdSggoydWb4tVQcBLfTSc81ylFwGL94Ws6ojC0eiMSDQdrcXeBycMIvdR7d5Q3TmmCrZGiWBhAsveg2kqHS5IlWWklKCrCYfIfznByxXFXiZhQ2Jvx6AbUQCJdzyIND4vrlOXn05wiGQ3oiiZn1mexDLuPKc6Gi7QBkhNQgyfC2cjHn2aO2bAJUL0lwzRyv64NPX8Ty9D3S0eR9TMET3b5"},
	//}
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", "CyY9AIlj6zRk3nVy")
	formData.Set("client_secret", "SmnvMxy3UiWZ89NNrHMtIQsg8C7AMljEsLoJzdJKYs1XrzMOBZmGYfIhKqeoOVDdu44ohHXmEcQGLaFYTw7rePSAfJcU8ToTmNg5yVnfMbgWsYRmcvAvU5sc0dK52oYaKwiT7UHMhM8sZie8OCxguh5ghj9rm0rdaAHZNNJ34bvRSS27Ixa5RYwuhJsG3EBoB5IEmxgT0WnhPLQLKqSX5ciokO1JiOUbfJswwypz6AKZhfdZCDRS3wpkNor9")
	formData.Set("refresh_token", "4OCe2FGdcneaF6cdSggoydWb4tVQcBLfTSc81ylFwGL94Ws6ojC0eiMSDQdrcXeBycMIvdR7d5Q3TmmCrZGiWBhAsveg2kqHS5IlWWklKCrCYfIfznByxXFXiZhQ2Jvx6AbUQCJdzyIND4vrlOXn05wiGQ3oiiZn1mexDLuPKc6Gi7QBkhNQgyfC2cjHn2aO2bAJUL0lwzRyv64NPX8Ty9D3S0eR9TMET3b5")

	var trVK TokenResponseVK

	//req, err = http.Post(url_full, "application/json", r)
	ctx := context.Background()
	fmt.Println("POST формы:", formData.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://ads.vk.com/api/v2/oauth2/token.json",
		bytes.NewReader([]byte(formData.Encode())),
	)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Выполняем запрос
	resp, err := client.Do(req)

	// Устанавливаем заголовки

	// Выполняем запрос

	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err), TokenResponseVK{}
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err), TokenResponseVK{}
	}

	err_response := json.Unmarshal(body, &trVK)

	if err_response != nil {
		return "", fmt.Errorf("ошибка сервера: %d, ответ: %s", resp, body), TokenResponseVK{}
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка сервера: %d, ответ: %s", resp.StatusCode, body), trVK
	}

	token_data, _ := json.Marshal(trVK)

	err_file := writeToVkConfig(vk_config_filename, string(token_data))
	if err_file != nil {

		fmt.Println("Ошибка записи файла конфигурации VK ", err_file.Error())
	}
	fmt.Println("ответ VK сервера:", string(resp.StatusCode), string(body), trVK)
	return string(body), nil, trVK
}

func GetUsersListsVK() (VK_Segments, error) {
	url := "https://ads.vk.com/api/v3/remarketing/users_lists.json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return VK_Segments{}, fmt.Errorf("failed to create request: %v", err)
	}

	// Устанавливаем заголовок Authorization
	req.Header.Set("Authorization", "Bearer Tm5HH1ZNhJqYDB0Eh3kaqkhQXKbJhPHs9jjX8Uy5pAi8CGDZX0vub4abt7TJzbhg2Ykygogbd4yzBEfh9Uy8Td6XFwDzNh0a9byKYBMUiJOka0k9tiBEFfxrIlyHFceqxkftXzDstmKAndRtw5vQoaYcUjSxNlPmRYJuA5IPF0wzaWLDd1FekH4hw9fqXAgkoH4XQmAhttJapeSJTUNm5ehvAlsG3794JBQeQw")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return VK_Segments{}, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return VK_Segments{}, fmt.Errorf("failed to read response: %v", err)
	}

	var vk_segments VK_Segments

	err = json.Unmarshal(body, &vk_segments)

	if err != nil {

		return VK_Segments{}, fmt.Errorf("failed to unmarshal response: %v", err)

	}

	return vk_segments, nil

}

func (vk VK_Segments) ToYandexSegments() (YaSegments, error) {

	var end_segments YaSegments

	if vk.Count > 0 {

		for i := range vk.Items {
			var temp YaSegment

			temp.ID = vk.Items[i].Id
			temp.Name = vk.Items[i].Name
			temp.ItemQuantity = vk.Items[i].IdsCount

			end_segments.Segments = append(end_segments.Segments, temp)

		}

	} else {

		return YaSegments{}, errors.New("no segments in the list")

	}

	return end_segments, nil

}
func getAccessTokenVK(apiURL, clientID, clientSecret string) (string, error, TokenResponseVK) {
	// Подготавливаем данные формы
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)

	var trVK TokenResponseVK

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %v", err), TokenResponseVK{}
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err), TokenResponseVK{}
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err), TokenResponseVK{}
	}

	err_response := json.Unmarshal(body, &trVK)

	if err_response != nil {
		return "", fmt.Errorf("ошибка сервера: %d, ответ: %s", resp.StatusCode, body), TokenResponseVK{}
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка сервера: %d, ответ: %s", resp.StatusCode, body), trVK
	}

	token_data, _ := json.Marshal(trVK)

	err_file := writeToVkConfig(vk_config_filename, string(token_data))
	if err_file != nil {

		fmt.Println("Ошибка записи файла конфигурации VK ", err_file.Error())
	}

	return string(body), nil, trVK
}
func refreshTokenVKCURL() (string, error) {
	// Подготовка команды curl
	cmd := exec.Command("curl",

		"-X", "POST",
		"-H", "Content-Type:application/x-www-form-urlencoded",
		"-d", "grant_type=refresh_token&client_id=CyY9AIlj6zRk3nVy&client_secret=SmnvMxy3UiWZ89NNrHMtIQsg8C7AMljEsLoJzdJKYs1XrzMOBZmGYfIhKqeoOVDdu44ohHXmEcQGLaFYTw7rePSAfJcU8ToTmNg5yVnfMbgWsYRmcvAvU5sc0dK52oYaKwiT7UHMhM8sZie8OCxguh5ghj9rm0rdaAHZNNJ34bvRSS27Ixa5RYwuhJsG3EBoB5IEmxgT0WnhPLQLKqSX5ciokO1JiOUbfJswwypz6AKZhfdZCDRS3wpkNor9&refresh_token=4OCe2FGdcneaF6cdSggoydWb4tVQcBLfTSc81ylFwGL94Ws6ojC0eiMSDQdrcXeBycMIvdR7d5Q3TmmCrZGiWBhAsveg2kqHS5IlWWklKCrCYfIfznByxXFXiZhQ2Jvx6AbUQCJdzyIND4vrlOXn05wiGQ3oiiZn1mexDLuPKc6Gi7QBkhNQgyfC2cjHn2aO2bAJUL0lwzRyv64NPX8Ty9D3S0eR9TMET3b5",
		"https://ads.vk.com/api/v2/oauth2/token.json",
	)

	// Захватываем stdout и stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Выполняем команду
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error executing curl: %v\nStderr: %s", err, stderr.String())
	}

	err_file := writeToVkConfig(vk_config_filename, out.String())
	if err_file != nil {

		fmt.Println("Ошибка записи файла конфигурации VK ", err_file.Error())
	}
	return out.String(), nil
}
func (c *TokenResponseVK) uploadUserListVK(full_url string, listType, listName, filePath string) error {
	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	// Создаем буфер для multipart запроса
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Добавляем поля формы
	_ = writer.WriteField("type", listType)
	_ = writer.WriteField("name", listName)

	// Создаем часть для файла
	part, err := writer.CreateFormFile("phones", filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания части файла: %v", err)
	}

	// Копируем содержимое файла
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("ошибка копирования файла: %v", err)
	}

	// Закрываем writer для финализации multipart запроса
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("ошибка закрытия writer: %v", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", full_url, &requestBody)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	fmt.Printf("Статус код: %d\n", resp.StatusCode)
	fmt.Printf("Ответ: %s\n", body)

	return nil
}
func writeToVkConfig(filename string, data string) error {
	// Получаем абсолютный путь к целевой директории
	targetDir := filepath.Join("..", "vk_config")

	// Проверяем существование директории
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("директория %s не существует", targetDir)
	}

	// Создаем полный путь к файлу
	filePath := filepath.Join(targetDir, filename)

	// Открываем файл для записи (создаем или перезаписываем)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %v", err)
	}
	defer file.Close()

	// Записываем данные
	_, err = file.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}

	return nil
}
