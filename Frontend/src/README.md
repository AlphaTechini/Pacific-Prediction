# Frontend/src

## Purpose

I use this folder for the actual frontend application code.

## Architectural Decisions And Tradeoffs

- I split shared UI/data code under `lib` and page ownership under `routes`.
- I keep backend transport details out of route components where possible.
- The tradeoff is more folders, but ownership stays much clearer as the app grows.

## Logic Tracking

- To find shared frontend code visit [lib/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/README.md).
- To find route ownership visit [routes/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/README.md).
- To find global app shell files visit [app.html](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/app.html), [app.css](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/app.css), and [app.d.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/app.d.ts).

## Component And Connection Map

- The shared frontend module layer can be found in [lib/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/README.md).
- The page route layer can be found in [routes/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/README.md).
