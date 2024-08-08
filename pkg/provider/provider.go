package provider

import (
	"encoding/json"
	"io"

	internal "github.com/daytonaio/daytona-provider-azure/internal"
	log_writers "github.com/daytonaio/daytona-provider-azure/internal/log"
	provider_types "github.com/daytonaio/daytona-provider-azure/pkg/types"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/provider"
	"github.com/daytonaio/daytona/pkg/provider/util"
	"github.com/daytonaio/daytona/pkg/workspace"
	"github.com/daytonaio/daytona/pkg/workspace/project"
)

type AzureProvider struct {
	BasePath           *string
	DaytonaDownloadUrl *string
	DaytonaVersion     *string
	ServerUrl          *string
	ApiUrl             *string
	LogsDir            *string
	ApiPort            *uint32
	ServerPort         *uint32
	NetworkKey         *string
	OwnProperty        string
}

func (p *AzureProvider) Initialize(req provider.InitializeProviderRequest) (*util.Empty, error) {
	p.OwnProperty = "my-own-property"

	p.BasePath = &req.BasePath
	p.DaytonaDownloadUrl = &req.DaytonaDownloadUrl
	p.DaytonaVersion = &req.DaytonaVersion
	p.ServerUrl = &req.ServerUrl
	p.ApiUrl = &req.ApiUrl
	p.LogsDir = &req.LogsDir
	p.ApiPort = &req.ApiPort
	p.ServerPort = &req.ServerPort
	p.NetworkKey = &req.NetworkKey

	return new(util.Empty), nil
}

func (p AzureProvider) GetInfo() (provider.ProviderInfo, error) {
	return provider.ProviderInfo{
		Name:    "azure-provider",
		Version: internal.Version,
	}, nil
}

func (p AzureProvider) GetTargetManifest() (*provider.ProviderTargetManifest, error) {
	return provider_types.GetTargetManifest(), nil
}

func (p AzureProvider) GetDefaultTargets() (*[]provider.ProviderTarget, error) {
	info, err := p.GetInfo()
	if err != nil {
		return nil, err
	}

	defaultTargets := []provider.ProviderTarget{
		{
			Name:         "default-target",
			ProviderInfo: info,
			Options:      "{\n\t\"Required String\": \"default-required-string\"\n}",
		},
	}
	return &defaultTargets, nil
}

func (p AzureProvider) CreateWorkspace(workspaceReq *provider.WorkspaceRequest) (*util.Empty, error) {
	logWriter := io.MultiWriter(&log_writers.InfoLogWriter{})
	if p.LogsDir != nil {
		loggerFactory := logs.NewLoggerFactory(*p.LogsDir)
		wsLogWriter := loggerFactory.CreateWorkspaceLogger(workspaceReq.Workspace.Id, logs.LogSourceProvider)
		logWriter = io.MultiWriter(&log_writers.InfoLogWriter{}, wsLogWriter)
		defer wsLogWriter.Close()
	}

	logWriter.Write([]byte("Workspace created\n"))

	return new(util.Empty), nil
}

func (p AzureProvider) StartWorkspace(workspaceReq *provider.WorkspaceRequest) (*util.Empty, error) {
	return new(util.Empty), nil
}

func (p AzureProvider) StopWorkspace(workspaceReq *provider.WorkspaceRequest) (*util.Empty, error) {
	return new(util.Empty), nil
}

func (p AzureProvider) DestroyWorkspace(workspaceReq *provider.WorkspaceRequest) (*util.Empty, error) {
	return new(util.Empty), nil
}

func (p AzureProvider) GetWorkspaceInfo(workspaceReq *provider.WorkspaceRequest) (*workspace.WorkspaceInfo, error) {
	providerMetadata, err := p.getWorkspaceMetadata(workspaceReq)
	if err != nil {
		return nil, err
	}

	workspaceInfo := &workspace.WorkspaceInfo{
		Name:             workspaceReq.Workspace.Name,
		ProviderMetadata: providerMetadata,
	}

	projectInfos := []*project.ProjectInfo{}
	for _, project := range workspaceReq.Workspace.Projects {
		projectInfo, err := p.GetProjectInfo(&provider.ProjectRequest{
			TargetOptions: workspaceReq.TargetOptions,
			Project:       project,
		})
		if err != nil {
			return nil, err
		}
		projectInfos = append(projectInfos, projectInfo)
	}
	workspaceInfo.Projects = projectInfos

	return workspaceInfo, nil
}

func (p AzureProvider) CreateProject(projectReq *provider.ProjectRequest) (*util.Empty, error) {
	logWriter := io.MultiWriter(&log_writers.InfoLogWriter{})
	if p.LogsDir != nil {
		loggerFactory := logs.NewLoggerFactory(*p.LogsDir)
		projectLogWriter := loggerFactory.CreateProjectLogger(projectReq.Project.WorkspaceId, projectReq.Project.Name, logs.LogSourceProvider)
		logWriter = io.MultiWriter(&log_writers.InfoLogWriter{}, projectLogWriter)
		defer projectLogWriter.Close()
	}

	logWriter.Write([]byte("Project created\n"))

	return new(util.Empty), nil
}

func (p AzureProvider) StartProject(projectReq *provider.ProjectRequest) (*util.Empty, error) {
	return new(util.Empty), nil
}

func (p AzureProvider) StopProject(projectReq *provider.ProjectRequest) (*util.Empty, error) {
	return new(util.Empty), nil
}

func (p AzureProvider) DestroyProject(projectReq *provider.ProjectRequest) (*util.Empty, error) {
	return new(util.Empty), nil
}

func (p AzureProvider) GetProjectInfo(projectReq *provider.ProjectRequest) (*project.ProjectInfo, error) {
	providerMetadata := provider_types.ProjectMetadata{
		Property: projectReq.Project.Name,
	}

	metadataString, err := json.Marshal(providerMetadata)
	if err != nil {
		return nil, err
	}

	projectInfo := &project.ProjectInfo{
		Name:             projectReq.Project.Name,
		IsRunning:        true,
		Created:          "Created at ...",
		ProviderMetadata: string(metadataString),
	}

	return projectInfo, nil
}

func (p AzureProvider) getWorkspaceMetadata(workspaceReq *provider.WorkspaceRequest) (string, error) {
	metadata := provider_types.WorkspaceMetadata{
		Property: workspaceReq.Workspace.Id,
	}

	jsonContent, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}

	return string(jsonContent), nil
}
