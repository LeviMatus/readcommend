import React, { useState, useEffect } from "react";
import Autocomplete from "@material-ui/lab/Autocomplete";
import TextField from "@material-ui/core/TextField";
import * as models from "../../models";

export default function Authors(props: {
  onChange: (authors: models.Author[]) => void;
}) {
  const [error, setError] = React.useState<models.Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [authors, setAuthors] = React.useState<models.Author[]>([]);

  useEffect(() => {
    fetch("http://localhost:5000/api/v1/authors")
      .then(res => res.json())
      .then(
        result => {
          setIsLoaded(true);
          setAuthors(result);
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
      options={authors}
      getOptionLabel={x => `${x.firstName} ${x.lastName}`}
      onChange={(_, items) => props.onChange(items)}
      renderInput={params => (
        <TextField
          {...params}
          variant="standard"
          label="Authors"
          placeholder="Type or select author(s)"
        />
      )}
    />
  );
}
