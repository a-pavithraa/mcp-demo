# MCP server for AWS services using Go

This is a sample implementation of a MCP server for AWS services built using the AWS SDK for Go. [mcp-go](https://github.com/mark3labs/mcp-go) project has been used as the MCP Go implementation.

This MCP server exposes the following tools for interacting with AWS services:

- **List DynamoDB Tables**: Retrieve a list of all DynamoDB tables in your AWS account.
- **Get DynamoDB Table Metadata**: Fetch metadata such as created date, size, and pricing model for DynamoDB tables.
- **List KMS Keys**: List all KMS keys along with metadata like created date, key description, and ID.
- **List S3 Buckets**: List all S3 buckets along with commonly used metadata like creation date, region, versioning status, etc.

## Sample Prompts

Here are some example prompts you can use with AI assistants that support MCP:

- "Show me a list of all my DynamoDB tables"
- "Can you get metadata for my DynamoDB tables?"
- "List all my KMS keys and their details"
- "What S3 buckets do I have and when were they created?"
- "Show me all my AWS storage resources and their metadata"
- "Which of my DynamoDB tables are using on-demand pricing?"
- "List my S3 buckets that have versioning enabled"

The AI will use the appropriate MCP tool to fetch the requested information directly from your AWS account.

## How to run

```bash
git clone https://github.com/apavithraa/mcp-demo
cd mcp-demo

go build -o awstools main.go
```

Configure the MCP server:

```bash
mkdir -p .vscode

# Define the content for mcp.json
MCP_JSON_CONTENT=$(cat <<EOF
{
  "servers": {
    "AWS Services MCP": {
      "type": "stdio",
      "command": "$(pwd)/awstools"
    }
  }
}
EOF
)

# Write the content to mcp.json
echo "$MCP_JSON_CONTENT" > .vscode/mcp.json
```

## AWS Authentication

- The MCP server uses the default AWS credential chain. Make sure you have configured your AWS credentials using one of the following methods:
  - AWS CLI (`aws configure`)
  - Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
  - IAM Roles for Amazon EC2 or ECS
  - AWS IAM Identity Center

You are good to go! Now spin up VS Code Insiders in Agent Mode, or any other MCP tool (like Claude Desktop) and try this out!

## Local dev/testing

Start with [MCP inspector](https://modelcontextprotocol.io/docs/tools/inspector) - `npx @modelcontextprotocol/inspector ./awstools`
