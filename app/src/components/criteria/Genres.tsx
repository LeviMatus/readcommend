import React, { useState, useEffect } from "react";
import Autocomplete from "@material-ui/lab/Autocomplete";
import TextField from "@material-ui/core/TextField";
import * as models from "../../models";

export default function Genres(props: {
  onChange: (genre: models.Genre[]) => void;
}) {
  const [error, setError] = useState<models.Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [genres, setGenres] = useState<models.Genre[]>([]);

  useEffect(() => {
    fetch("http://localhost:5000/api/v1/genres")
      .then(res => res.json())
      .then(
        result => {
          setIsLoaded(true);
          setGenres(result);
        },
        error => {
          setIsLoaded(true);
          setError(error);
        }
      );
  }, []);

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  if (!isLoaded) {
    return <div>Loading...</div>;
  }

  return (
    <Autocomplete
      multiple
      options={genres!}
      getOptionLabel={x => x.title}
      onChange={(_, items) => props.onChange(items)}
      renderInput={params => (
        <TextField
          {...params}
          variant="standard"
          label="Genres"
          placeholder="Type or select genre(s)"
        />
      )}
    />
  );
}
