package manager

import (
	"errors"
	"strings"
	"testing"

	"github.com/WolkenOps/sekeep-api/internal/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

//#########################################
// Mocks for AWS SDK SSM
// We can unit test functions that uses
// SSM client
//#########################################
type mockSSMClient struct {
	ssmiface.SSMAPI
}

var mockedPassword = ssm.Parameter{
	Name:  aws.String("/sekeep/work/mail/myhandler"),
	Value: aws.String("password.lol"),
}

var mockedParametersMetadata = []*ssm.ParameterMetadata{&ssm.ParameterMetadata{Name: mockedPassword.Name}}

func (m *mockSSMClient) DeleteParameter(input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	if *input.Name == *mockedPassword.Name {
		return &ssm.DeleteParameterOutput{}, nil
	}
	return nil, errors.New("ParameterNotFound")
}

func (m *mockSSMClient) DescribeParameters(input *ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error) {
	output := ssm.DescribeParametersOutput{
		Parameters: mockedParametersMetadata,
	}
	return &output, nil
}

func (m *mockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	if *input.Name == *mockedPassword.Name {
		return &ssm.GetParameterOutput{
			Parameter: &mockedPassword,
		}, nil
	}

	return nil, errors.New("ParameterNotFound")
}

func (m *mockSSMClient) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	if *input.Name == *mockedPassword.Name {
		return &ssm.PutParameterOutput{
			Tier:    aws.String("ok"),
			Version: aws.Int64(111111),
		}, nil
	} else if strings.Contains(*input.Name, "exists") {
		return nil, errors.New("ParameterAlreadyExists")
	}
	return nil, errors.New("Uncatched Error")
}

// #########################################
// Tests starts here
//##########################################
func TestCreateOrUpdate(t *testing.T) {
	testPassword := model.Password{Name: "/work/mail/myhandler", Value: "password.lol"}
	ssmClient = &mockSSMClient{}

	//Test if creates a new password
	err := CreateOrUpdate(testPassword)
	if err != nil {
		t.Errorf("CreateOrUpdate failed with error, %s", err.Message)
	}

	testPassword.Name = "/work/exists"

	//Test if parameter already exists
	err = CreateOrUpdate(testPassword)
	if err.StatusCode != 409 {
		t.Errorf("CreateOrUpdate failed, expected %d, got %d", 409, err.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	testPassword := model.Password{Name: "/work/mail/myhandler", Value: "password.lol"}
	ssmClient = &mockSSMClient{}

	//Test if creates a new password
	err := Delete(testPassword)
	if err != nil {
		t.Errorf("Delete failed with error, %s", err.Message)
	}

	testPassword.Name = "/not/found"

	//Test if parameter already exists
	err = Delete(testPassword)
	if err.StatusCode != 404 {
		t.Errorf("Delete failed, expected %d, got %d", 409, err.StatusCode)
	}
}

func TestRead(t *testing.T) {
	testPassword := model.Password{Name: "/work/mail/myhandler"}
	ssmClient = &mockSSMClient{}

	//Test if Read returns a valid password
	password, _ := Read(testPassword)
	if *mockedPassword.Value == password {
		t.Log("Read returned a valid password")
	} else {
		t.Errorf("Read failed, expected %s, got %s", password, *mockedPassword.Value)
	}

	testPassword.Name = "/not/exist"

	//Read returns error if password is not found
	_, err := Read(testPassword)
	if err.StatusCode == 404 {
		t.Log("Read returned 404")
	} else {
		t.Errorf("Read failed, expected %d, got %d", 404, err.StatusCode)
	}
}

func TestList(t *testing.T) {
	testPassword := model.Password{Name: "work"}
	ssmClient = &mockSSMClient{}
	passwords, _ := List(testPassword)
	if len(passwords) == 1 {
		t.Log("List returned one element")
	} else {
		t.Errorf("List failed, expected %d, got %d", 1, len(passwords))
	}
}
