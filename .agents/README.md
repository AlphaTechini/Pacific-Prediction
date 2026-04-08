# .agents

## Purpose

I use this folder to keep persistent project guidance that should survive across planning and implementation turns.

## Architectural Decisions And Tradeoffs

- I keep durable constraints here instead of scattering them across chat-only context.
- I use `GUIDE.md` as the short operational memory file because it is faster to scan than the full architecture document.
- The tradeoff is duplication, but I prefer that over losing important decisions during implementation.

## Logic Tracking

- To find confirmed project constraints visit [GUIDE.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/GUIDE.md).
- To find the high-level architecture that those constraints came from visit [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The persistent decision memory can be found in [GUIDE.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/GUIDE.md).
- The repo-level structure map can be found in [structure.md](file:///C:/Hackathons/Pacific%20Prediction/structure.md).
