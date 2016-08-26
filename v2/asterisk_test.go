package ari

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsteriskInfo(t *testing.T) {
	assert := assert.New(t)
	//This only tests to make sure there are no errors.
	info, err := DefaultClient.GetAsteriskInfo("")
	if err != nil {
		t.Error("Got error getting asterisk info:", err.Error())
	}
	// Check buildinfo
	assert.NotEmpty(info.BuildInfo.Date)
	// Check ConfigInfo
	//assert.NotEmpty(info.ConfigInfo.Name) // Does not appear to be implemented
	assert.NotEmpty(info.ConfigInfo.DefaultLanguage)
	// Check SetId (execution user/group)
	//assert.NotEmpty(info.ConfigInfo.SetId.Group)  // Does not appear to be implemented
	//assert.NotEmpty(info.ConfigInfo.SetId.User)   // Does not appear to be implemented
	// Check StatusInfo
	assert.NotEmpty(info.StatusInfo.StartupTime)
	assert.NotEmpty(info.StatusInfo.LastReloadTime)
	// Check SystemInfo
	assert.NotEmpty(info.SystemInfo.EntityID)
	assert.NotEmpty(info.SystemInfo.Version)

}

func TestAsteriskGlobalVariables(t *testing.T) {
	assert := assert.New(t)

	err := DefaultClient.SetAsteriskVariable("testVar", "testVal")
	if err != nil {
		t.Error("Got error setting asterisk global variable:", err.Error())
	}

	val, err := DefaultClient.GetAsteriskVariable("testVar")
	if err != nil {
		t.Error("Got error getting asterisk global variable:", err.Error())
	}
	assert.Equal("testVal", val)
}
