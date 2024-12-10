# Groupie Tracker

## Description

`groupie-tracker` is a user-friendly web app that shows artist-related information provided by four APIs:
* **[`artists`](https://groupietrackers.herokuapp.com/api/artists)** is a catalog of the artist image, name, members, creation date, first album, etc.
* **[`locations`](https://groupietrackers.herokuapp.com/api/locations)** list concert locations and dates.
* **[`dates`](https://groupietrackers.herokuapp.com/api/dates)** list only the concert dates.
* **[`relation`](https://groupietrackers.herokuapp.com/api/relation)** links the location and date information.

The information from the APIs is displayed on the webpage inside blocks that are responsive and interactive. Clicking the artist image or name opens a collapsible panel that shows more details, such as a link to the **Artist Page**. Clicking the **Artist Page** link triggers an event/action (`GET` request). This is just one example of the HTTP request-response client-and-server communication implemented in this project.

## Setup and Usage

### Prerequisites

* Go: version 1.19 or higher
* A compatible web browser
* An active internet connection for API data

### Usage

1. Clone the repository:
```bash
git clone https://github.com/OthmaneAfilali/Groupie-Tracker.git
cd groupie-tracker
```

2. Execute the program:
```bash
go run .
```

3. Open your browser and visit [`http://localhost:8080/`](http://localhost:8080/) or [`http://localhost:8080/groupie-tracker`](http://localhost:8080/groupie-tracker).

## Features

### Functional

* The backend is written in Go, using only standard packages.
* Errors have been handled, and the following are examples:
	* **`400 Bad Request`: Failed to fetch data from API.** This happens if someone navigates to a link that specifies the wrong artist ID.
	* **`404 Status Not Found`: URL not found.** This happens if someone types in the wrong URL path.
	* **`500 Internal Server Error`: Error parsing API JSON.** This happens if there is an error processing the API.
	* **`502 Bad Gateway`: Failed to fetch data from API.** This happens, for example, if the API link being accessed is incorrect.

### Responsive

Here are some examples of the CSS features we implemented to make the site responsive:

* Flex display is used to ensure artist boxes fluidly grows or shrinks to fit the space available.
* Relative sizing for fonts (rem) and viewport height and width (vh and vw, respectively) are implemented to respond to different screen settings.

### Deployment

The tracker has been deployed at [`https://groupie-tracker-e3mz.onrender.com`](https://groupie-tracker-e3mz.onrender.com)

## [The Groupies](http://localhost:8080/groupie-tracker/about)

* [Markus Amberla](https://github.com/MarkusYPA/MarkusYPA)
* [Othmane Afilali](https://github.com/OthmaneAfilali)
* [Jedi Reston](https://github.com/jeeeeedi)