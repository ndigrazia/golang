package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const jsonFile = "Uploadfile.json"

const nameTemplateFile = "PatchPurpose.tmpl"

const path = "templates/" + nameTemplateFile

var templates map[string]*template.Template

//func to substract inside template
var funcMap = template.FuncMap{
	"minus": func(a, b int) int {
		return a - b
	},
}

//FourthPlatformConnectionData represents the token from an HTTP response.
type FourthPlatformConnectionData struct {
	AccessToken, TokenType, Scope, Purpose string
	ExpiresIn                              float64
	manager                                *FourthPlatformManagerData
	contextBody                            map[string]string
	contextHeader                          map[string]string
}

//FourthPlatformManagerData represents the connection data.
type FourthPlatformManagerData struct {
	Username, Password, URIAuth, URIAdmin string
}

//ResourceNotFoundError represents no found errors.
type ResourceNotFoundError struct {
	message, uri string
}

//ResponseHTTPError represents request errors.
type ResponseHTTPError struct {
	description, uri, response string
	code                       int
}

//InternalError represents request errors.
type InternalError struct {
	message string
	err     error
}

// PiScopeData struct which contains an id,
// a title and a description of Pi Scope
type PiScopeData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

//NewPurposeData which contains an id, a title, a description, and Pi Scopes of Purpose
//Use to create a new purpose
type NewPurposeData struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Level       string        `json:"level"`
	PiScopes    []PiScopeData `json:"pi_scopes"`
}

//UpdatePurposeData struct which contains an id,
// an pi-scopes of Purpose
type UpdatePurposeData struct {
	ID       string        `json:"id_purpose"`
	PiScopes []PiScopeData `json:"pi_scopes"`
}

// FourthPlatformFileData struct which contains an array of PiScope, an purpose id and a Swagger Url
// Struct loaded from the file strings.json
type FourthPlatformFileData struct {
	PurposeID    string        `json:"id_purpose"`
	PurposeDesc  string        `json:"desc_purpose"`
	PurposeLevel string        `json:"level_purpose"`
	SwaggerURL   string        `json:"swagger_url"`
	PiScopes     []PiScopeData `json:"pi_scopes"`
}

//APISpecificationData struct
type APISpecificationData struct {
	Version string
	Prefix  string
}

func initContext(managers map[string]*FourthPlatformManagerData) {
	log.Output(2, "LOADING environment..")

	var keys [2]string

	keys[0] = "dev"
	managers[keys[0]] = &FourthPlatformManagerData{
		//Username: "baikal-api",
		//Password: "MzEyZjRjOTZkZDI0ZDA0NTAwMGVlNmEw",
		URIAuth:  "https://auth.ar-dev.baikalplatform.com/token",
		URIAdmin: "https://api.ar-dev.baikalplatform.com/admin/v2",
	}

	keys[1] = "pre"
	managers[keys[1]] = &FourthPlatformManagerData{
		//Username: "baikal-api",
		//Password: "QJukMUUVNq",
		URIAuth:  "https://auth.ar-pre.baikalplatform.com/token",
		URIAdmin: "https://api.ar-pre.baikalplatform.com/admin/v2",
	}

	log.Output(2, fmt.Sprint("Environment LOADED: ", keys))

	//Template
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["PatchPurpose"] = template.Must(template.New("patch").Funcs(funcMap).ParseFiles(path))

	log.Output(2, fmt.Sprint("Template loaded successfully: ", path))
}

func loadContext(managers map[string]*FourthPlatformManagerData, env string, user string, pass string) {
	//Setting user and pass
	if value, ok := managers[strings.ToLower(env)]; !ok {
		log.Output(2, fmt.Sprint("The enviroment ", env, " doesn't exist."))
		os.Exit(1)
	} else {
		log.Output(2, fmt.Sprint("WORKING in ", env))
		value.Password = pass
		value.Username = user
	}
}

//Error returns an error.
func (e *ResourceNotFoundError) Error() string {
	return e.message + "-" + e.uri
}

//Error returns an error.
func (e *InternalError) Error() string {
	var r strings.Builder
	r.WriteString(e.message)

	if e.err != nil {
		r.WriteString(" - ")
		r.WriteString(e.err.Error())
	}

	return r.String()
}

//Error returns an error.
func (e *ResponseHTTPError) Error() string {
	return e.description + "-" + e.uri + "-" + e.response + "-" + strconv.Itoa(e.code)
}

//GetConnection returns a connection.
func (a *FourthPlatformManagerData) GetConnection(body map[string]string, header map[string]string) (*FourthPlatformConnectionData, error) {
	var result *FourthPlatformConnectionData

	data := url.Values{}

	if body != nil {
		for k, v := range body {
			data.Set(k, v)
		}
	}

	input, code, err := sendRequest("POST", a.Username, a.Password, a.URIAuth, strings.NewReader(data.Encode()), header)
	if err != nil {
		return result, err
	}

	resp := make(map[string]interface{})
	err = json.Unmarshal(input, &resp)
	if err != nil {
		return result, err
	}

	if code != 200 {
		return result, &ResponseHTTPError{"http: Cannot get token for target", a.URIAuth, string(input), code}
	}

	result = &FourthPlatformConnectionData{
		AccessToken:   resp["access_token"].(string),
		ExpiresIn:     resp["expires_in"].(float64),
		Purpose:       resp["purpose"].(string),
		Scope:         resp["scope"].(string),
		TokenType:     resp["token_type"].(string),
		manager:       a,
		contextBody:   body,
		contextHeader: header,
	}

	return result, nil
}

//ExistsResource determines if an resource exists.
func (a *FourthPlatformConnectionData) ExistsResource(resource string) (bool, error) {
	var r strings.Builder
	r.WriteString(a.manager.URIAdmin)
	r.WriteString(resource)

	var header = make(map[string]string)

	var buf strings.Builder
	buf.WriteString(a.TokenType)
	buf.WriteString(" ")
	buf.WriteString(a.AccessToken)
	header["Authorization"] = buf.String()

	resp, code, err := sendRequest("GET", "", "", r.String(), nil, header)
	if err != nil {
		return false, err
	}
	if code != 200 && code != 404 {
		return false, &ResponseHTTPError{"Cannot identify the resource", r.String(), string(resp), code}
	}

	return code == 200, nil
}

//CreateResource creates a resource.
func (a *FourthPlatformConnectionData) CreateResource(resource string, body []byte) error {
	var r strings.Builder
	r.WriteString(a.manager.URIAdmin)
	r.WriteString(resource)

	var header = make(map[string]string)

	var buf strings.Builder
	buf.WriteString(a.TokenType)
	buf.WriteString(" ")
	buf.WriteString(a.AccessToken)
	header["Authorization"] = buf.String()
	header["Content-Type"] = "application/json"

	input, code, err := sendRequest("POST", "", "", r.String(), bytes.NewBuffer(body), header)
	if err != nil {
		return err
	}

	if code != 201 && code != 202 {
		return &ResponseHTTPError{"http: Cannot create the resource", r.String(), string(input), code}
	}

	return nil
}

//PatchResource updates a resource.
func (a *FourthPlatformConnectionData) PatchResource(resource string, body []byte) error {
	var r strings.Builder
	r.WriteString(a.manager.URIAdmin)
	r.WriteString(resource)

	var header = make(map[string]string)

	var buf strings.Builder
	buf.WriteString(a.TokenType)
	buf.WriteString(" ")
	buf.WriteString(a.AccessToken)
	header["Authorization"] = buf.String()
	header["Content-Type"] = "application/json"

	input, code, err := sendRequest("PATCH", "", "", r.String(), bytes.NewBuffer(body), header)
	if err != nil {
		return err
	}

	if code != 202 && code != 201 {
		return &ResponseHTTPError{"http: Cannot patch the resource", r.String(), string(input), code}
	}

	return nil
}

//Reconnect reconnects
func (a *FourthPlatformConnectionData) Reconnect() error {
	c, err := a.manager.GetConnection(a.contextBody, a.contextHeader)
	if err != nil {
		return err
	}

	a.AccessToken = c.AccessToken
	a.ExpiresIn = c.ExpiresIn
	a.Purpose = c.Purpose
	a.Scope = c.Scope
	a.TokenType = c.TokenType

	return nil
}

//GetResource returns an resource.
func (a *FourthPlatformConnectionData) GetResource(resource string) ([]byte, error) {
	var r strings.Builder
	r.WriteString(a.manager.URIAdmin)
	r.WriteString(resource)

	var header = make(map[string]string)

	var buf strings.Builder
	buf.WriteString(a.TokenType)
	buf.WriteString(" ")
	buf.WriteString(a.AccessToken)
	header["Authorization"] = buf.String()

	resp, code, err := sendRequest("GET", "", "", r.String(), nil, header)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		if code == 404 {
			return nil, &ResourceNotFoundError{"Cannot identify the resource", r.String()}
		}
		return nil, &ResponseHTTPError{"Error to get the resource", r.String(), string(resp), code}
	}

	return resp, nil
}

//DeleteResource remove an resource.
func (a *FourthPlatformConnectionData) DeleteResource(resource string) error {
	var r strings.Builder
	r.WriteString(a.manager.URIAdmin)
	r.WriteString(resource)

	var header = make(map[string]string)

	var buf strings.Builder
	buf.WriteString(a.TokenType)
	buf.WriteString(" ")
	buf.WriteString(a.AccessToken)
	header["Authorization"] = buf.String()

	resp, code, err := sendRequest("DELETE", "", "", r.String(), nil, header)
	if err != nil {
		return err
	}
	if code != 200 && code != 202 {
		return &ResponseHTTPError{"Cannot delete the resource.", r.String(), string(resp), code}
	}

	return nil
}

func sendRequest(method string, username string, passwd string, uri string, body io.Reader, header map[string]string) ([]byte, int, error) {
	var bodyText []byte = nil

	client := &http.Client{}
	req, _ := http.NewRequest(method, uri, body)
	if username != "" && passwd != "" {
		req.SetBasicAuth(username, passwd)
	}

	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return bodyText, 0, err
	}

	defer resp.Body.Close()

	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bodyText, 0, err
	}

	return bodyText, resp.StatusCode, nil
}

func loadJSONFile(fourthPlatformFile *FourthPlatformFileData, name string) error {
	//Open our jsonFile
	jsonFile, err := os.Open(name)
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	log.Output(2, fmt.Sprint("Successfully opened file: ", name))

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, fourthPlatformFile)
	if err != nil {
		return err
	}

	log.Output(2, fmt.Sprint("Successfully parsed file: ", name))

	return nil
}

func createPiScope(conn *FourthPlatformConnectionData, piscope PiScopeData) {
	//Exists Resource
	var r strings.Builder
	r.WriteString("/pi-scopes/")
	r.WriteString(piscope.ID)

	log.Output(2, fmt.Sprint("loading PI-Scope: ", piscope.ID))

	exists, err := conn.ExistsResource(r.String())
	if err != nil {
		log.Panicln("Could not load pi-scope: ", r.String(), err)
	}

	//Create resource
	if !exists {
		log.Output(2, fmt.Sprint("PI-Scope doesn´t exist: ", piscope.ID))
		value, err := json.Marshal(piscope)
		if err != nil {
			log.Panicln("Could not marshal pi-scope: ", r.String(), err)
		}

		err = conn.CreateResource("/pi-scopes", value)
		if err != nil {
			log.Panicln("Could not create pi-scope: ", r.String(), err)
		}

		log.Output(2, fmt.Sprint("PI-Scope created successfully.", piscope.ID))
	} else {
		log.Output(2, fmt.Sprint("PI-Scope exists: ", piscope.ID))
	}

	log.Output(2, fmt.Sprint("PI-Scope done: ", piscope.ID))
}

func getJSONPurpose(purpose *UpdatePurposeData) ([]byte, error) {
	var buf bytes.Buffer

	//Applies a parsed template to the slice of Note objects
	err := templates["PatchPurpose"].ExecuteTemplate(&buf, nameTemplateFile, purpose)
	if err != nil {
		return nil, err
	}

	readBuf, err := ioutil.ReadAll(&buf)
	if err != nil {
		return nil, err
	}

	return readBuf, nil
}

func patchPurpose(conn *FourthPlatformConnectionData, purpose *UpdatePurposeData) error {
	readBuf, err := getJSONPurpose(purpose)
	if err != nil {
		return err
	}

	err = conn.PatchResource("/purposes", readBuf)
	if err != nil {
		return err
	}
	return nil
}

func getPiScopesInPurpose(conn *FourthPlatformConnectionData, purposeID string) ([]PiScopeData, error) {
	var r strings.Builder
	r.WriteString("/purposes/")
	r.WriteString(purposeID)

	resp, err := conn.GetResource(r.String())
	if err != nil {
		return nil, err
	}

	type LangData struct {
		LangID string `json:"lang_id"`
		Value  string `json:"value"`
	}

	var purpose struct {
		ID               string     `json:"id"`
		Title            []LangData `json:"title"`
		Description      []LangData `json:"description"`
		ShortDescription []LangData `json:"short_description"`
		OptIn            []LangData `json:"opt_in"`
		OptOut           []LangData `json:"opt_out"`
		Level            string     `json:"level"`
		PiScopes         []struct {
			ID          string     `json:"id"`
			Title       []LangData `json:"title"`
			Description []LangData `json:"description"`
		} `json:"pi_scopes"`
	}

	err = json.Unmarshal(resp, &purpose)
	if err != nil {
		return nil, err
	}

	PiScopeAux := make([]PiScopeData, len(purpose.PiScopes))

	for i := 0; i < len(purpose.PiScopes); i++ {
		PiScopeAux[i] = PiScopeData{
			ID:          purpose.PiScopes[i].ID,
			Title:       purpose.PiScopes[i].Title[0].Value,
			Description: purpose.PiScopes[i].Description[0].Value,
		}
	}

	return PiScopeAux, nil
}

func union(piscope1 []PiScopeData, piscope2 []PiScopeData) []PiScopeData {
	if piscope1 == nil {
		return piscope2
	}

	if piscope2 == nil {
		return piscope1
	}

	union := append(piscope1, piscope2...)
	seen := make(map[string]struct{}, len(union))

	j := 0
	for _, piscope := range union {
		if _, ok := seen[piscope.ID]; ok {
			continue
		}
		seen[piscope.ID] = struct{}{}
		union[j] = piscope
		j++
	}

	return union[:j]
}

func createBackup(purpose *UpdatePurposeData) error {
	log.Output(2, "Backuping...")

	var fileLog strings.Builder
	fileLog.WriteString("patchbackupfile")
	fileLog.WriteString("_")
	fileLog.WriteString(time.Now().Format("20060102150405"))
	fileLog.WriteString(".log")

	f, err := os.OpenFile(fileLog.String(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	//JSON Purpose
	log.Output(2, "Creating json for backuping...")
	jsonPurposeForBackup, err := getJSONPurpose(purpose)
	if err != nil {
		return err
	}

	log.Output(2, "json for backuping DONE.")

	log.Output(2, "writing to backup file...")

	_, err = f.WriteString(string(jsonPurposeForBackup))
	if err != nil {
		return err
	}

	log.Output(2, "Backup DONE.")

	return nil
}

func createPurpose(conn *FourthPlatformConnectionData, purpose *NewPurposeData) error {
	log.Output(2, fmt.Sprint("Creating purpose: ", purpose.ID))

	value, err := json.Marshal(purpose)
	if err != nil {
		return err
	}

	err = conn.CreateResource("/purposes", value)
	if err != nil {
		return err
	}

	log.Output(2, fmt.Sprint("Purpose CREATED."))

	return nil
}

func doPatchWithPurpose(conn *FourthPlatformConnectionData, piScopeCurrent []PiScopeData, fourthPlatformFile *FourthPlatformFileData) {
	log.Output(2, "Merging pi-scopes...")
	mergedPiScope := union(piScopeCurrent, fourthPlatformFile.PiScopes)
	log.Output(2, "Pi-scopes MERGED.")

	purposeUpdate := &UpdatePurposeData{
		ID:       fourthPlatformFile.PurposeID,
		PiScopes: mergedPiScope,
	}
	purposeCurrent := &UpdatePurposeData{
		ID:       fourthPlatformFile.PurposeID,
		PiScopes: piScopeCurrent,
	}

	err := createBackup(purposeCurrent)
	if err != nil {
		log.Panicln("Could not write in the backup file.", err)
	}

	log.Output(2, "PATCHING...")
	err = patchPurpose(conn, purposeUpdate)
	if err != nil {
		log.Output(2, fmt.Sprintln("Could not patch the purpose: ", fourthPlatformFile.PurposeID, err))
		if purposeCurrent != nil {
			log.Output(2, "Patching purpose backuped...")
			err = patchPurpose(conn, purposeCurrent)
			if err != nil {
				log.Output(2, fmt.Sprintln("WARNING - Could not patch purpose backuped: ", fourthPlatformFile.PurposeID, err))
				log.Panicln("WARNING: Please, patch manually the purpose with the last backup file.")
			}
			log.Output(2, "Patch DONE.")
		}
	} else {
		log.Output(2, "Patch DONE.")
	}
}

func doPatchWithoutPurpose(conn *FourthPlatformConnectionData, fourthPlatformFile *FourthPlatformFileData) {
	purpose := &NewPurposeData{
		ID:          fourthPlatformFile.PurposeID,
		Title:       fourthPlatformFile.PurposeID,
		Description: fourthPlatformFile.PurposeDesc,
		Level:       fourthPlatformFile.PurposeLevel,
		PiScopes:    fourthPlatformFile.PiScopes,
	}
	err := createPurpose(conn, purpose)
	if err != nil {
		log.Panicln(2, fmt.Sprintln("WARNING - Could not create purpose: ", fourthPlatformFile.PurposeID, err))
	}
}

func loadSpecificationFrom(spec *APISpecificationData, swaggerURL string) error {
	resp, code, err := sendRequest("GET", "", "", swaggerURL, nil, nil)
	if err != nil {
		return err
	}

	if code != 200 {
		return &ResponseHTTPError{"Could not load 4ta Plataform API Specifications:", swaggerURL, string(resp), code}
	}

	var apiSpecificationData struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
		APIPrefix string `json:"x-fp-apiPrefix"`
	}

	err = json.Unmarshal(resp, &apiSpecificationData)
	if err != nil {
		return &InternalError{"Could not parse 4ta Plataform API.", err}
	}

	//Getting API's version and name to load to enviroment
	s := strings.Split(apiSpecificationData.Info.Version, ".")
	if len(s) == 0 {
		return &InternalError{"Minimum version not found.", err}
	}

	spec.Version = strings.Split(apiSpecificationData.Info.Version, ".")[0]
	spec.Prefix = apiSpecificationData.APIPrefix

	return nil
}

func uploadAPI(conn *FourthPlatformConnectionData, spec *APISpecificationData, fourthPlatformFile *FourthPlatformFileData) error {
	var r strings.Builder
	r.WriteString("/apis")
	r.WriteString(spec.Prefix)
	r.WriteString("/v")
	r.WriteString(spec.Version)

	log.Output(2, fmt.Sprint("READING 4ta Plataform Enviroment to find API: ", r.String()))
	exists, err := conn.ExistsResource(r.String())
	if err != nil {
		return err
	}

	if !exists {
		log.Output(2, "API doesn´t exist. ")
		err = createAPI(conn, fourthPlatformFile)
		if err != nil {
			return err
		}
	} else {
		log.Output(2, "API exists.")
		log.Output(2, "DELETING API.")
		err = conn.DeleteResource(r.String())
		if err != nil {
			return err
		}
		log.Output(2, "API DELETED.")
		err = createAPI(conn, fourthPlatformFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func createAPI(conn *FourthPlatformConnectionData, fourthPlatformFile *FourthPlatformFileData) error {
	log.Output(2, fmt.Sprint("Creating API: ", fourthPlatformFile.SwaggerURL))

	var api struct {
		SwaggerURL string `json:"swagger_url"`
	}

	api.SwaggerURL = fourthPlatformFile.SwaggerURL

	value, err := json.Marshal(api)
	if err != nil {
		return err
	}

	err = conn.CreateResource("/apis", value)
	if err != nil {
		return err
	}

	log.Output(2, fmt.Sprint("API CREATED."))

	return nil
}

func main() {
	var enviroment, username, password string

	flag.StringVar(&enviroment, "env", "", "4th Plataform enviroment")
	flag.StringVar(&enviroment, "e", "", "4th Plataform enviroment (shorthand)")

	flag.StringVar(&username, "user", "", "4th Plataform username")
	flag.StringVar(&username, "u", "", "4th Plataform username (shorthand)")

	flag.StringVar(&password, "pass", "", "4th Plataform password")
	flag.StringVar(&password, "p", "", "4th Plataform password (shorthand)")

	flag.Parse()

	if enviroment == "" || username == "" || password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var managers map[string]*FourthPlatformManagerData = make(map[string]*FourthPlatformManagerData)

	//Context
	log.Output(2, "GETTING Context to 4ta Plataform...")
	initContext(managers)
	loadContext(managers, enviroment, username, password)
	log.Output(2, "Context GOTTEN.")

	//Open our json configuration File
	log.Output(2, fmt.Sprint("UPLOADING the file ", jsonFile, " to 4ta Plataform..."))
	var fourthPlatformFile *FourthPlatformFileData = &FourthPlatformFileData{}
	err := loadJSONFile(fourthPlatformFile, jsonFile)
	if err != nil {
		log.Panicln("Could not load file.", jsonFile, err)
	}

	//Load API specification
	var spec *APISpecificationData = nil

	if fourthPlatformFile.SwaggerURL != "" {
		spec = &APISpecificationData{}
		log.Output(2, fmt.Sprint("LOADING 4ta Plataform API Specifications: ", fourthPlatformFile.SwaggerURL))
		err = loadSpecificationFrom(spec, fourthPlatformFile.SwaggerURL)
		if err != nil {
			log.Panicln("Could not load ta Plataform API Specifications: ", fourthPlatformFile.SwaggerURL, err)
		}
		log.Output(2, fmt.Sprint("4ta Plataform API Specifications LOADED. Url: ", spec.Prefix, ", Versión: ", spec.Version))
	} else {
		log.Output(2, "Swagger URL NO DEFINED.")
	}

	//Getting connection from 4ta Plataform
	body := make(map[string]string)
	body["grant_type"] = "client_credentials"
	body["scope"] = "admin:pi-scopes:create admin:pi-scopes:read admin:pi-scopes:delete admin:purposes:read admin:purposes:update admin:purposes:create admin:apis:read admin:apis:create admin:apis:delete"

	log.Output(2, fmt.Sprint("Body loaded: ", body))

	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"

	log.Output(2, fmt.Sprint("Header loaded: ", header))

	log.Output(2, "GETTING connection from 4ta Plataform....")

	conn, err := managers[enviroment].GetConnection(body, header)
	if err != nil {
		log.Panicln("Error Connection", jsonFile, err)
	}
	log.Output(2, "Connection DONE.")

	//PiScopes
	lPiScopes := len(fourthPlatformFile.PiScopes)
	if lPiScopes > 0 {
		log.Output(2, "LOADING PI-Scopes...")
		for i := 0; i < lPiScopes; i++ {
			createPiScope(conn, fourthPlatformFile.PiScopes[i])
		}
		log.Output(2, "PI-Scopes DONE.")
	} else {
		log.Output(2, "PI-Scopes NO DEFINED.")
	}

	//API
	if spec != nil {
		log.Output(2, "UPLOADING API...")
		err = uploadAPI(conn, spec, fourthPlatformFile)
		if err != nil {
			log.Panicln("Could not upload API.", err)
		}
		log.Output(2, "API UPLOADED.")
	}

	//Purpose
	if fourthPlatformFile.PurposeID != "" {
		log.Output(2, fmt.Sprint("LOADING purpose: ", fourthPlatformFile.PurposeID))
		piScopeCurrent, err := getPiScopesInPurpose(conn, fourthPlatformFile.PurposeID)
		if err != nil {
			switch err.(type) {
			case *ResourceNotFoundError:
				log.Output(2, fmt.Sprintln("Don't exist purpose: ", fourthPlatformFile.PurposeID))
				//Patching
				doPatchWithoutPurpose(conn, fourthPlatformFile)
			default:
				log.Panicln("Could not load purpose: ", fourthPlatformFile.PurposeID, err)
			}
		} else {
			//Patching
			log.Output(2, "Purpose LOADED.")
			doPatchWithPurpose(conn, piScopeCurrent, fourthPlatformFile)
		}
	} else {
		log.Output(2, "Purpose NO DEFINED.")
	}

	log.Output(2, "DONE.")
}
