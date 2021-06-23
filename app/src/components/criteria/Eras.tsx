import React, { useState, useEffect } from "react";
import * as models from "../../models";
import { MenuItem, Select } from "@material-ui/core";

export default function Eras(props: { onChange: (era: models.Era) => void }) {
  const [error, setError] = useState<models.Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [eras, setEras] = useState<models.Era[]>([]);

  useEffect(() => {
    fetch("http://localhost:5000/api/v1/eras")
      .then(res => res.json())
      .then(
        result => {
          setIsLoaded(true);
          setEras(result);
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
    <Select
      onChange={e => props.onChange(eras[Number(e.target.value)])}
      fullWidth={true}
      defaultValue={0}
    >
      {eras.map((x, i) => (
        <MenuItem key={x.id} value={i}>
          {x.title}
        </MenuItem>
      ))}
    </Select>
  );
}
