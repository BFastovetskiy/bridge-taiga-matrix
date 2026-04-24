# bridge-taiga-matrix

A lightweight Go service that monitors deadlines in [Taiga](https://taiga.io) projects and sends notifications to [Matrix](https://matrix.org) rooms.

## How it works

On each run the service:
1. Authenticates with the Taiga API using login/password.
2. Iterates over the configured projects.
3. Fetches all open user stories and tasks for each project.
4. For every item checks:
   - no deadline set → sends a warning;
   - deadline is in the past → sends an overdue alert;
   - deadline is within `daysUntilDeadline` days → sends a "days left" reminder.
5. Sends the message to the project's Matrix room and, optionally, to a general room.

Designed to be run as a scheduled job (cron, Task Scheduler, systemd timer, etc.).

## Requirements

- Go 1.22+
- Access to a Taiga instance (self-hosted or cloud)
- A Matrix account with a valid access token and room membership

## Configuration

Copy `settings-example.json` to `settings.json` and fill in the values.

```json
{
    "taigaBaseURL"  : "https://taiga.example.com",
    "taigaUsername" : "your-username",
    "taigaPassword" : "your-password",
    "taigaProjects" : [
        {
            "name": "my-project-slug",
            "matrixProjectRoomID": "!roomid:matrix.example.com"
        }
    ],

    "matrixServer" : "https://matrix.example.com",
    "matrixToken"  : "syt_your_matrix_token",
    "duplicateToGeneralGroup": true,
    "generalRoomId": "!generalroomid:matrix.example.com",

    "InsecureSkipVerify": false,

    "language": "en",

    "daysUntilDeadline": 15
}
```

### Parameters

| Parameter | Type | Description |
|---|---|---|
| `taigaBaseURL` | string | Base URL of your Taiga instance |
| `taigaUsername` | string | Taiga login |
| `taigaPassword` | string | Taiga password |
| `taigaProjects` | array | List of projects to monitor (see below) |
| `taigaProjects[].name` | string | Project slug (visible in the Taiga URL) |
| `taigaProjects[].matrixProjectRoomID` | string | Matrix room ID for this project |
| `matrixServer` | string | Base URL of your Matrix homeserver |
| `matrixToken` | string | Matrix access token |
| `duplicateToGeneralGroup` | bool | Also send all notifications to `generalRoomId` |
| `generalRoomId` | string | Matrix room ID for the general channel |
| `InsecureSkipVerify` | bool | Skip TLS certificate verification (use `false` in production) |
| `language` | string | Notification language: `en` or `ru` (empty string defaults to system locale) |
| `daysUntilDeadline` | int | Send a reminder when deadline is this many days away or fewer |

## Build & run

```bash
go build -o bridge-taiga-matrix .
./bridge-taiga-matrix -config settings.json
```

The `-config` flag is optional; `settings.json` in the current directory is used by default.

## Localization

Notification templates are stored in `locales/<lang>.json`. Add a new file to support additional languages and set `"language"` in the config accordingly.

## Scheduling example (Linux cron)

Run every morning at 09:00:

```cron
0 9 * * * /opt/bridge-taiga-matrix/bridge-taiga-matrix -config /opt/bridge-taiga-matrix/settings.json
```
