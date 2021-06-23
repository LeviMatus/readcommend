import React from "react";
import image from "../images/banner.jpg";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
  fullWidthImage: {
    background: `black url(${image}) center no-repeat`,
    height: "390px",
    position: "relative",
  },
  titleShader: {
    backgroundColor: "#0000007F",
    width: "100%",
    height: "140px",
    position: "absolute",
    bottom: "0",
    left: "0",
  },
  titleBox: {
    maxWidth: "960px",
    marginLeft: "auto",
    marginRight: "auto",
  },
  titleText: {
    margin: "0px",
    paddingLeft: "30px",
    fontFamily: "Caveat",
    fontSize: "60pt",
    color: "white",
    verticalAlign: "top",
  },
  subtitleText: {
    margin: "0px",
    paddingLeft: "30px",
    fontFamily: "Caveat",
    fontSize: "20pt",
    color: "white",
    verticalAlign: "top",
    position: "relative",
    top: "-20px",
  },
}));

export default function Banner() {
  const classes = useStyles();
  return (
    <div className={classes.fullWidthImage}>
      <div className={classes.titleShader}>
        <div className={classes.titleBox}>
          <h1 className={classes.titleText}>Readcommend</h1>
          <h2 className={classes.subtitleText}>
            Book recommendations better than from your best friend
          </h2>
        </div>
      </div>
    </div>
  );
}
