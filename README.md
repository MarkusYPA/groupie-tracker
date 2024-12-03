# Groupie Tracker

## Overview

Groupie Tracker is a dynamic and interactive web application that showcases information about bands and artists using a provided API. The website provides detailed data visualizations, including information about artists, their members, concert dates, and locations. It supports server-client interaction to deliver a seamless and engaging user experience.

This project is built using **Go** and adheres to robust development practices, ensuring smooth operation and error handling.

---

## Features

### Core Functionalities:
- **Artist Profiles**: View artist details, including name, image, members, year of activity, and first album date.
- **Concert Info**: Explore past and upcoming concert locations and dates.
- **Data Visualizations**: User-friendly presentation of data through cards, tables, lists, and graphics.
- **Event-Driven Actions**: Features dynamic client-server interaction, such as submitting artist names via forms and retrieving specific data in real time.
  
### Highlights:
- **Responsive Design**: Styled for optimal viewing across devices.
- **Error Handling**: Displays informative error pages to ensure smooth navigation.
- **Testable Codebase**: Includes unit tests for verifying application functionality.

---

## Project Structure

### Files and Folders:
```plaintext
/groupie-tracker
│
├── main.go
├── internal/
│   ├── api/
│   │   └── api.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── shared/
│   │   └── shared.go
│   └── utils/
│       └── utils.go
├── templates/
│   ├── index.html
│   ├── about.html
│   ├── bio.html
│   └── error.html
├── assets/
│   ├── style.css
│   └── script.js
└── .gitignore
```
### File Descriptions

#### `main.go`
**Purpose:** Entry point of the application.  
**Responsibilities:**  
- Initializes templates and sets up HTTP routes.  
- Fetches and validates data from the API.  
- Starts the web server.  

**Key Functions:**  
- **`init()`**: Initializes the HTML templates.  
- **`main()`**: Fetches data, validates it, sets up routes, and launches the server.  

#### `internal/api/api.go`
**Purpose:** Handles communication with the API and data processing.  
**Responsibilities:**  
- Fetches and parses API data.  
- Validates the integrity and completeness of the data.  

**Key Functions:**  
- **`FetchAllData`**: Aggregates data from the API into a `PageData` struct.  
- **`fetchAndUnmarshal`**: Fetches JSON data from the API and unmarshals it into a provided interface.  
- **`FetchApi`**: Makes HTTP GET requests and returns the response body.  
- **`ValidateData`**: Ensures that fetched data is accurate and complete.  

#### `internal/handlers/handlers.go`
**Purpose:** Manages server-side logic for routing and rendering pages.  
**Responsibilities:**  
- Handles HTTP requests and responses for different routes.  
- Processes form submissions and manages user interactions.  

**Key Functions:**  
- **`HomeHandler`**: Renders the home page template with data.  
- **`AboutHandler`**: Renders the about page template with relevant information.  
- **`BioHandler`**: Manages GET and POST requests for artist details, delegating to `handleBioGet` or `handleBioPost`.  
- **`handleBioGet`**: Fetches artist data for GET requests and renders the bio page.  
- **`handleBioPost`**: Validates form input for POST requests and redirects to the bio page.  

#### `internal/shared/types.go`
**Purpose:** Defines shared types and data structures.  
**Key Types:**  
- **`Index`**: Represents index data from the API.  
- **`People`**: Represents an artist or band.  
- **`Location`**: Represents concert locations.  
- **`Date`**: Represents concert dates.  
- **`Relation`**: Links artists, dates, and locations.  
- **`PageData`**: Aggregates all fetched data types.  
- **`Page`**: Represents a page with a header and message.  

#### `internal/shared/templates.go`
**Purpose:** Stores shared template variables and references.  
**Key Variables:**  
- **`Data`**: Holds the fetched and aggregated API data.  
- **`IndexTmpl`**: Home page template reference.  
- **`AboutTmpl`**: About page template reference.  
- **`BioTmpl`**: Bio page template reference.  
- **`ErrTmpl`**: Error page template reference.  

#### `internal/utils/utils.go`
**Purpose:** Utility functions for error handling and logging.  
**Key Functions:**  
- **`ErrorHandler`**: Renders appropriate error pages based on HTTP status codes.  
- **`ErrorCheck`**: Logs errors if present and returns a status indicator.  
- **`LogError`**: Logs error messages with custom details.  

#### `templates/`
**Purpose:** Contains HTML templates for rendering different application pages.  
**Files:**  
- **`index.html`**: Home page template.  
- **`about.html`**: About page template.  
- **`bio.html`**: Artist details page template.  
- **`error.html`**: Error handling and messaging template.  

#### `assets/`
**Purpose:** Holds static files to enhance user experience.  
**Files:**  
- **`style.css`**: Provides custom styling for the application.  
- **`script.js`**: Contains client-side logic, including form submissions and dynamic interactions.  

#### `.gitignore`
**Purpose:** Specifies files and directories to ignore in the repository.  
**Contents:**  
- Common files like `.DS_Store` or other system-generated files.  
- Temporary or sensitive files that should not be committed to version control.  

---

## Setup and Usage

### Prerequisites
- **Go**: Version 1.19 or higher.
- A compatible web browser.
- An active internet connection for API data.

### Running the Application
1. Clone the repository:
   ```bash
   git clone https://github.com/OthmaneAfilali/Groupie-Tracker.git
   cd groupie-tracker
2. Execute the program:
   ```bash
   go run .
3. Open your browser and visit:
    ```bash
    http://localhost:8080/groupie-tracker

## API Integration

This application integrates with the [Groupie Tracker API](https://groupietrackers.herokuapp.com/api). The API is divided into four endpoints:

1. **`artists`**: Provides general information about bands and artists.
2. **`locations`**: Contains concert locations.
3. **`dates`**: Provides concert schedules.
4. **`relation`**: Links all the above data for a cohesive view.

---

## Contributing

Contributions are welcome! If you'd like to contribute:
1. Fork the repository.
2. Create a new branch (`feature/my-feature`).
3. Submit a pull request with detailed changes.

Please follow good coding practices and ensure your contributions include relevant tests.
