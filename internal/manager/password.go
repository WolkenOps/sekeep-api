package manager

import (
	"fmt"
	"strings"

	"github.com/WolkenOps/sekeep-api/internal/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	log "github.com/sirupsen/logrus"
)

var (
	ssmClient    ssmiface.SSMAPI
	seekepPrefix = "/sekeep"
)

type PasswordError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *PasswordError) Error() string {
	return fmt.Sprintf("Password failed with message %s, status code is %d", e.Message, e.StatusCode)
}

func init() {
	session := session.Must(session.NewSession())
	ssmClient = ssm.New(session)
}

func CreateOrUpdate(password model.Password) *PasswordError {
	log.Infof("createOrUpdate started on password %s", password.Name)
	_, err := ssmClient.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(seekepPrefix + password.Name),
		Value:     aws.String(password.Value),
		Type:      aws.String("SecureString"),
		Overwrite: aws.Bool(password.Overwrite)})
	if err != nil {
		if strings.Contains(err.Error(), "ParameterNotFound") {
			log.Warnf("Password %s not found", password.Name)
			return &PasswordError{"Parameter Not Found", 404, err}
		}
		if strings.Contains(err.Error(), "ParameterAlreadyExists") {
			log.Warnf("Password %s already exists", password.Name)
			return &PasswordError{"Password already exists", 409, err}
		}
		log.Errorf("createOrUpdate failed on password %s, %s", password.Name, err)
		return &PasswordError{fmt.Sprintf("Internal Error: %s", err), 500, err}
	}
	return nil
}

func Delete(password model.Password) *PasswordError {
	log.Infof("delete started on password %s", password.Name)
	_, err := ssmClient.DeleteParameter(&ssm.DeleteParameterInput{
		Name: aws.String(seekepPrefix + password.Name),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ParameterNotFound") {
			log.Warnf("Password %s not found", password.Name)
			return &PasswordError{"Parameter Not Found", 404, err}
		}
		log.Errorf("delete failed on password %s, %s", password.Name, err)
		return &PasswordError{fmt.Sprintf("Internal Error: %s", err), 500, err}
	}
	return nil
}

func Read(password model.Password) (string, *PasswordError) {
	log.Infof("read started on password %s", password.Name)
	parameter, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(seekepPrefix + password.Name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ParameterNotFound") {
			log.Warnf("Password %s not found", password.Name)
			return "", &PasswordError{"Parameter Not Found", 404, err}
		}

		log.Errorf("read failed on password %s, %s", password.Name, err)
		return "", &PasswordError{fmt.Sprintf("Internal Error: %s", err), 500, err}
	}
	return *parameter.Parameter.Value, nil
}

func List(password model.Password) ([]model.Password, *PasswordError) {
	if len(password.Name) == 0 {
		password.Name = "/"
	}
	log.Infof("list started on path %s", password.Name)
	filter := ssm.ParameterStringFilter{
		Key:    aws.String("Name"),
		Option: aws.String("Contains"),
		Values: aws.StringSlice([]string{string(seekepPrefix + password.Name)}),
	}

	parameters, err := ssmClient.DescribeParameters(&ssm.DescribeParametersInput{
		ParameterFilters: []*ssm.ParameterStringFilter{&filter},
	})

	if err != nil {
		log.Errorf("list failed on path %s, %s", password.Name, err)
		return []model.Password{}, &PasswordError{fmt.Sprintf("Internal Error: %s", err), 500, err}
	}

	log.Debugf("Passwords found: %d", len(parameters.Parameters))
	var passwords = []model.Password{}

	for _, p := range parameters.Parameters {
		name := *p.Name
		newName := strings.Join(remove(strings.Split(name, "/"), 1), "/")
		passwords = append(passwords, model.Password{Name: newName})
	}

	return passwords, nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
