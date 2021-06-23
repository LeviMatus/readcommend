import React, { useState, useEffect } from "react";
import * as models from "../../models";
import { MenuItem, Select } from "@material-ui/core";

export default function Sizes(props: {
  onChange: (size: models.Size) => void;
}) {
  const [error, setError] = useState<models.Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [sizes, setSizes] = useState<models.Size[]>([]);

  useEffect(() => {
    fetch("http://localhost:5000/api/v1/sizes")
      .then(res => res.json())
      .then(
        result => {
          setIsLoaded(true);
          setSizes(result);
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
      onChange={e => props.onChange(sizes[Number(e.target.value)])}
      fullWidth={true}
      defaultValue={0}
    >
      {sizes.map((x, i) => (
        <MenuItem key={x.id} value={i}>
          {x.title}
        </MenuItem>
      ))}
    </Select>
  );
}
