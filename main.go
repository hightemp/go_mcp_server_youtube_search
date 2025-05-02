package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/hightemp/youtube-search-api-go/pkg/youtubesearchapi"
	"github.com/hightemp/youtube-transcript-api-go/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport string
	var host string
	var port string
	flag.StringVar(&transport, "t", "sse", "Transport type (stdio or sse)")
	flag.StringVar(&host, "h", "0.0.0.0", "Host of sse server")
	flag.StringVar(&port, "p", "8889", "Port of sse server")
	flag.Parse()

	mcpServer := server.NewMCPServer(
		"go_mcp_server_mdurl",
		"1.0.0",
	)

	youtubeSearchTool := mcp.NewTool("youtube_search",
		mcp.WithDescription("Find videos on youtube"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("text query"),
		),
	)

	mcpServer.AddTool(youtubeSearchTool, youtubeSearchHandler)

	youtubeGetVideoInfoTool := mcp.NewTool("youtube_get_video_info",
		mcp.WithDescription("Get video details"),
		mcp.WithString("videoID",
			mcp.Required(),
			mcp.Description("videoID, for example aBcDeFgHiJk"),
		),
	)

	mcpServer.AddTool(youtubeGetVideoInfoTool, youtubeGetVideoInfoHandler)

	youtubeGetSubtitlesTool := mcp.NewTool("youtube_get_subtitles",
		mcp.WithDescription("Get youtube video subtitles"),
		mcp.WithString("videoID",
			mcp.Required(),
			mcp.Description("videoID, for example aBcDeFgHiJk"),
		),
	)

	mcpServer.AddTool(youtubeGetSubtitlesTool, youtubeGetSubtitles)

	if transport == "sse" {
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL(fmt.Sprintf("http://localhost:%s", port)))
		log.Printf("SSE server listening on %s:%s URL: http://127.0.0.1:%s/sse", host, port, port)
		if err := sseServer.Start(fmt.Sprintf("%s:%s", host, port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

func sprintJSON(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	return string(jsonData)
}

func youtubeSearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, errors.New("query must be a string")
	}

	searchResult, err := youtubesearchapi.GetData(query, true, 5, nil)
	if err != nil {
		log.Fatalf("Error getting data: %v", err)
	}

	result := sprintJSON(searchResult)

	return mcp.NewToolResultText(result), nil
}

func youtubeGetVideoInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	videoID, ok := request.Params.Arguments["videoID"].(string)
	if !ok {
		return nil, errors.New("videoID must be a string")
	}

	videoDetails, err := youtubesearchapi.GetVideoDetails(videoID)
	if err != nil {
		return nil, fmt.Errorf("Error getting video details: %v", err)
	}

	result := sprintJSON(videoDetails)

	return mcp.NewToolResultText(result), nil
}

func youtubeGetSubtitles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	videoID, ok := request.Params.Arguments["videoID"].(string)
	if !ok {
		return nil, errors.New("videoID must be a string")
	}

	youtubeAPI := api.NewYouTubeTranscriptApi()
	languages := []string{"en", "ru"}

	transcript, err := youtubeAPI.GetTranscript(videoID, languages)
	if err != nil {
		return nil, fmt.Errorf("Error getting transcript: %v", err)
	}

	result := ""
	for _, entry := range transcript.Entries {
		result += fmt.Sprintf("[%.2f - %.2f]: %s\n", entry.Start, entry.Start+entry.Duration, entry.Text)
	}

	return mcp.NewToolResultText(result), nil
}
