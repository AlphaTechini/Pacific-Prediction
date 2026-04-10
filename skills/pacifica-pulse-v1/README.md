# Pacifica Pulse v1 Skill

## Purpose

I use this folder to package a public OpenClaw skill that teaches an agent the current Pacifica Pulse v1 boundaries, architecture defaults, and settlement rules.

## Architectural Decisions And Tradeoffs

- I publish this skill under the product name `pacifica-pulse-v1` so the ClawHub slug stays stable and easy to discover.
- I keep the first release in a single `SKILL.md` file because the current guidance fits comfortably in one place and stays easier to review before publishing.
- I focus the skill on Pacifica Pulse v1 and similar read-only Pacifica prediction products instead of generic trading systems so the instructions stay opinionated and actually useful.
- I avoid extra runtime requirements or installer metadata because this skill is documentation-driven and should load cleanly on a default OpenClaw setup.

## Logic Tracking

- To find the publishable skill instructions visit [SKILL.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/SKILL.md).
- To find the skill folder index visit [../README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/README.md).
- To find the source architecture decisions visit [../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- To find the confirmed product constraints visit [../../.agents/GUIDE.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/GUIDE.md).

## Component And Connection Map

- The OpenClaw skill entrypoint can be found in [SKILL.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/SKILL.md).
- The repository structure map can be found in [../../structure.md](file:///C:/Hackathons/Pacific%20Prediction/structure.md).
