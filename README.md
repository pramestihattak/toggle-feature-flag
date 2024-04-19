# Toggle Feature Flag

This is a proof of concept on how we can manage feature flag in database and making changes without restarting the service. Please note that we *DO NOT* retrieve data from db every single time.
We only making db connection when
1. Starting the service
2. State changes in db

# Concept
1. The state of feature flag is stored in database.
2. Whenever the service is started, it loads the state from database.
3. The state then stored in feature flag handler with hashmap.
4. Feature flag handler will be injected to the services.
5. We can use a getter method from handler to get all feature flags state `GetFeatureFlags()`
6. There will be one endpoint that can be used to manipulate feature flag state (enabled/disabled). *(TODO)*
7. Whenever there's a change in feature flag table, postgres will notify the handler.
8. We do have listener in feature flag handler, so it will know that there's a state change in feature flag state.
9. The handler will reconstruct the hashmap and will use the up to date state from db.
10. Now the changes will be reflected to the service as well.

# How to run
1. Create postgres db instance with docker `docker run --name postgres-db -e POSTGRES_PASSWORD=secret -p 5438:5432 -d postgres`
2. Dump `data.sql` manually
3. Try hit `localhost:3000` by default it should return `enabled` 
4. Try to change db state in `toggle_features` table (This should be done via API endpoint, but since this is just POC, let's just change it directly from db).
5. Redo step #3, it should return `disabled` now.
