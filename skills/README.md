# Skills

## Purpose

I use this folder to keep publishable OpenClaw skill packages that can be installed outside this repository without depending on private local paths.

## Architectural Decisions And Tradeoffs

- I keep public skill payloads self-contained so another user can install them from ClawHub and get useful behavior immediately.
- I start with one compact skill instead of a larger multi-file package because the first release should optimize for clarity and low maintenance overhead.
- I keep the public skill aligned with the repo's real product boundary so the published instructions do not drift into unsupported wallet, trading, or AI-summary features.

## Logic Tracking

- To find the current public OpenClaw skill package visit [pacifica-pulse-v1/SKILL.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/SKILL.md).
- To find the skill packaging notes visit [pacifica-pulse-v1/README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/README.md).
- To find the repo-wide structure map visit [../structure.md](file:///C:/Hackathons/Pacific%20Prediction/structure.md).
- To find the confirmed project constraints visit [../.agents/GUIDE.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/GUIDE.md).

## Component And Connection Map

- The public OpenClaw skill package can be found in [pacifica-pulse-v1/SKILL.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/SKILL.md).
- The skill folder ownership notes can be found in [pacifica-pulse-v1/README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/README.md).
