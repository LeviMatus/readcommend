import React from "react";
import { MenuItem, Select } from "@material-ui/core";

export default function Limit(props: { onChange: (limit: number) => void }) {
  const limits = [5, 10, 20, 50];
  return (
    <Select
      onChange={e => props.onChange(Number(e.target.value))}
      fullWidth={true}
      defaultValue={10}
    >
      {limits.map(x => (
        <MenuItem key={x} value={x}>
          {x}
        </MenuItem>
      ))}
    </Select>
  );
}
