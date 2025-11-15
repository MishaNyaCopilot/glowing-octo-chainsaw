# AniStream Project Improvements

Here is a list of planned improvements to make "AniStream" a more robust and feature-rich platform.

### Ⅰ. Backend & Infrastructure Enhancements (Production-Readiness)

1.  **Structured Logging:**
    *   **Problem:** The current use of the standard `log` package is not ideal for production. It lacks structure, making it difficult to parse, filter, and analyze logs effectively.
    *   **Suggestion:** Replace the standard `log` package with a structured logging library like **`Zerolog`** or **`Zap`**. This will provide leveled, JSON-formatted logs that are much easier to manage in a production environment.

2.  **Configuration Management:**
    *   **Problem:** Configuration is loaded directly from environment variables in `config.go`. This is a good start, but it's not very flexible.
    *   **Suggestion:** Integrate a library like **`Viper`**. It allows you to manage configuration from multiple sources (environment variables, config files like YAML/JSON, remote K/V stores) with a clear precedence order.

3.  **Database Migrations:**
    *   **Problem:** The project relies on GORM's `AutoMigrate`, which is convenient for development but risky for production as it doesn't handle complex schema changes or rollbacks well.
    *   **Suggestion:** Implement a dedicated database migration tool like **`golang-migrate/migrate`**. This will give you version-controlled, explicit SQL migration files, ensuring safe and predictable schema evolution.

4.  **Request Validation:**
    *   **Problem:** There is no explicit validation for incoming API request bodies.
    *   **Suggestion:** Add validation to your Gin handlers. You can integrate the **`go-playground/validator`** library with GORM models using struct tags to automatically validate incoming JSON payloads.

5.  **Raw Video File Cleanup:**
    *   **Problem:** The project plan mentions the need for a policy to delete raw video files after transcoding, but this is not yet implemented.
    *   **Suggestion:** Create a new RabbitMQ message queue (e.g., `video_cleanup_queue`). After a video is successfully transcoded, the worker can publish a message with the raw file's path. A new, simple worker (or a goroutine in the existing worker) can consume from this queue and delete the corresponding file from MinIO.

### Ⅱ. Frontend Enhancements (UI/UX & Polish)

1.  **UI Overhaul with a Component Library:**
    *   **Problem:** The current UI is functional but basic. It uses a mix of custom CSS and some unutilized Shadcn UI setup.
    *   **Suggestion:** Fully commit to using **Shadcn UI**. Systematically replace the custom components and styles in `App.css` with components from the library (`Card`, `Button`, `Input`, `Dialog` for modals, etc.). This will create a more modern, consistent, and professional-looking interface.

2.  **Improved State Management:**
    *   **Problem:** State is managed with `useState` and passed down through props, which can become cumbersome as the application grows.
    *   **Suggestion:** Introduce a client-side state management library like **`Zustand`** or **`Redux Toolkit`**. This will centralize application state (like the list of animes, search results, etc.), making it easier to manage and access from any component.

3.  **User Feedback (Notifications):**
    *   **Problem:** The admin panel performs actions (upload, delete) without giving clear visual feedback to the user.
    *   **Suggestion:** Add a "toast" notification system. Libraries like **`react-hot-toast`** are easy to integrate and can be used to show success or error messages after API calls (e.g., "Episode deleted successfully").

4.  **Loading States & Skeletons:**
    *   **Problem:** Pages show a simple "Loading..." text, which can feel abrupt.
    *   **Suggestion:** Implement **skeleton loaders**. While data is being fetched for the home page or anime detail page, display placeholder cards that mimic the final layout. This provides a much smoother perceived loading experience.

### Ⅲ. New Features

1.  **User Authentication & Profiles:**
    *   **Suggestion:** This is the most significant feature addition. Implement a user authentication system (e.g., using JWTs). This will unlock many new possibilities:
        *   **Watch History:** Track which episodes a user has watched.
        *   **Favorites/Watchlist:** Allow users to save anime to a personal list.
        *   **Continue Watching:** Show users where they left off.

2.  **Advanced Search & Filtering:**
    *   **Suggestion:** Enhance the search functionality:
        *   Create a dedicated backend endpoint (`/api/search`) that performs database queries.
        *   Add filters to the UI to allow users to search by **genre**, **release year**, etc.
        *   Implement **pagination** on the backend and frontend to efficiently handle large lists of anime.

3.  **Genre System:**
    *   **Suggestion:** Fully implement the `genres` field mentioned in the project plan.
        *   Create `genres` and `anime_genres` tables in the database.
        *   Add API endpoints to manage genres and associate them with anime.
        *   Create a page or section on the frontend where users can browse anime by genre.

4.  **Subtitle Support:**
    *   **Suggestion:** Add support for video subtitles.
        *   Update the backend to allow uploading subtitle files (e.g., `.vtt`) associated with an episode.
        *   Modify the `VideoPlayer` component on the frontend to load and display these subtitles.
