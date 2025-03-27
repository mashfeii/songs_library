# Song Library API

The Song Library API is a RESTful service for managing a song collection, supporting CRUD operations and paginated verse retrieval.

## API Endpoints

- **GET /songs**: Retrieve a paginated list of songs.
- **GET /songs/{id}/verses**: Retrieve a paginated list of verses for a song.
- **POST /songs**: Create a new song.
- **PUT /songs/{id}**: Update a song by ID.
- **DELETE /songs/{id}**: Delete a song by ID.

## Running the Application

- **Locally**:

```sh
docker run -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=songdb -p 5432:5432 postgres
make build && ./bin/songs_library
```

- **Docker**:

```sh
docker-compose -f docker/docker-compose.yml up --build
```
