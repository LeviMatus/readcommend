import React from "react";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
  copyright: {
    fontFamily: "Caveat",
    fontSize: "14pt",
    textAlign: "center",
    color: "#d0d0d0",
    marginTop: "24px",
    marginBottom: "80px",
  },
}));

export default function Copyright() {
  const classes = useStyles();
  return (
    <div className={classes.copyright}>
      Copyright (C)2021 nesto, inc. — All Rights Reserved
    </div>
  );
}
