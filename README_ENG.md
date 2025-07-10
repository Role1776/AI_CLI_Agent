# AI CLI Agent

AI CLI Agent is a powerful command-line assistant powered by Large Language Models (LLMs) designed to simplify your terminal workflow. It understands natural language requests, generates shell commands, and can autonomously propose, execute, and analyze system commands — no more googling bash or PowerShell syntax!

---

![AI CLI Agent Screenshot](https://raw.githubusercontent.com/Role1776/AI_CLI_Agent/main/photo_2025-07-10_15-18-47.jpg)


⚠️ Important Note: This tool is intended for educational purposes only. It is designed to simplify terminal workflows and help users understand command execution, but it should not be used in production environments or critical projects. It may contain bugs and is not optimized for performance.

## Features

- **Dual Modes:**
  - **Agent Mode (default):** Describe tasks in plain English (e.g., "create a folder named 'test' and navigate into it"), and the agent generates the command, asks for your confirmation, executes it, then summarizes the result.
  - **Chat Mode:** Prefix your query with `!` to get direct AI answers without command execution.

- **Safe Execution:** Always asks for confirmation before running any command (unless auto-execution mode is enabled).

- **Automatic Retries:** If a command fails, the agent analyzes the error, suggests fixes, and retries automatically. The number of retries is configurable.

- **Cross-Platform Support:** Works with PowerShell on Windows and bash on Linux/macOS.

- **Result Analysis:** After execution, the agent summarizes what happened or analyzes errors to provide insightful feedback.

- **User-Friendly Interface:** Includes colored output, loading spinners, and well-formatted code blocks for easy reading.

- **Session Management:** Maintains conversation history during your session. Use `/clear` to reset history anytime.

---

## Installation & Setup

### Prerequisites

- Go 1.20+ installed on your system.

### Installation Steps

1. **Clone the repository:**

    ```sh
    git clone https://github.com/your-username/cli_agent.git
    cd cli_agent
    ```

2. **Create a `.env` configuration file:**

    Example `.env` contents:

    ```env
    API_URL=https://openrouter.ai/api/v1/chat/completions
    API_TOKEN=your_api_openrouter_key_here
    MODEL=qwen/qwen3-30b-a3b:free
    RETRIES=3        # Recommended: 2-5
    TIMEOUT_CLIENT=30 # HTTP client timeout in seconds
    ```

3. **Build the application:**

    ```sh
    go build -o cli-agent ./cmd/cli_agent
    ```

4. **Run the CLI Agent:**

    ```sh
    ./cli-agent
    ```

---

## Usage

### Agent Mode (default)

Simply type your task in natural language:

```sh
> find all text files in the current directory
 ```
The agent will generate the command, ask for your approval, execute it, and provide a summary.

### Chat Mode
Prefix your message with ! to get an AI answer without executing commands:

```sh
> !what are goroutines in Go?
 ```

## Special Commands
### Command	Description
 ```sh
- /exit	        #Exit the program 
- /clear	#Clear the current session history 
- /auto-true	#Enable auto-execution mode (no confirmation) 
- /auto-false	#Disable auto-execution mode (default) 
 ```
    

    
   
