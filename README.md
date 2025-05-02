# go_mcp_server_youtube_search

A simple MCP (Model Context Protocol) server that provides tools for searching and retrieving information from YouTube videos. This server can be used with AI assistants that support the MCP protocol.

Server url: `http://127.0.0.1:8889/sse`

## Tools

- **youtube_search** - Search for YouTube videos using a text query
- **youtube_get_video_info** - Get detailed information about a video by its ID
- **youtube_get_subtitles** - Get video subtitles in English or Russian

## Usage

The server can operate in two modes: stdio and sse (Server-Sent Events). By default, it uses the sse mode.

```bash
go_mcp_server_youtube_search -t sse -h 0.0.0.0 -p 8889
# or
go_mcp_server_youtube_search -t stdio
```