# AI CLI Agent

AI CLI Agent is a powerful command-line assistant leveraging Large Language Models (LLMs) to simplify your terminal workflow. It can answer questions in a chat-like manner and autonomously generate, propose for execution, and analyze system commands.

Forget about googling PowerShell or bash syntax—just describe your task in natural language, and the agent will handle it for you!

---

## Features

- **Dual-Mode Operation:**
  - **Agent Mode:** Describe a task (e.g., "create a folder named 'test' and navigate into it"), and the agent will generate and execute the required command.
  - **Chat Mode:** Prefix your query with `!` to get a direct answer from the AI without executing any commands.

- **Safe Execution:** The agent asks for your confirmation before running any generated commands.

- **Automatic Retries:** If a command fails, the agent attempts to fix it and suggests a new version. The number of retries is configurable.

- **Cross-Platform:** Supports PowerShell on Windows and bash on Linux/macOS.

- **Result Analysis:** After executing a command successfully, the AI summarizes what was done; if an error occurs, it analyzes the cause.

- **User-Friendly Interface:** Includes colored output, loading spinners, and well-formatted blocks for easy reading.

- **Session Management:** Maintains conversation history within a session, which can be cleared with the `/clear` command.

---

## Installation and Setup

### Prerequisites

- Go (version 1.20 or higher) installed.

### Steps

1. **Clone the repository:**

   ```sh
   git clone https://github.com/your-username/cli_agent.git
   cd cli_agent
   
2. Create configuration file:

Create .env

Fill in .env:

API_URL: LLM provider API endpoint.

API_TOKEN: Your API key.

MODEL: AI model to use.

RETRIES: Recommended value between 2 and 5.

TIMEOUT_CLIENT: http timeout

Build the application:

    ```sh
    go build -o cli-agent ./cmd/cli_agent

Agent Mode (default)
Type your request in natural language:

    ```sh
    > find all text files in the current directory

The agent will generate the command, ask for confirmation, execute it, and provide a summary.

Chat Mode
Prefix your query with !:

    ```sh
    > !what are goroutines in Go?
## Special Commands

/exit	Exit the program
/clear	Clear current session history
/auto-true	Enable auto-execution mode (no confirmation)
/auto-false	Disable auto-execution mode (default)    
    

    
   
