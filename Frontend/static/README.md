# Frontend/static

## Purpose

I use this folder for static frontend assets that SvelteKit serves directly.

## Architectural Decisions And Tradeoffs

- I keep truly static public assets here instead of routing them through component code.
- This keeps the app shell simpler, but these files should stay small and intentional.
- I use `favicon.png` as both the browser favicon and the navigation brand image.

## Logic Tracking

- To find the frontend favicon and nav brand image visit [favicon.png](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/favicon.png).
- To find the current static robots configuration visit [robots.txt](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/robots.txt).

## Component And Connection Map

- The browser favicon and nav brand image can be found in [favicon.png](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/favicon.png).
- The static public asset boundary can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/README.md).
