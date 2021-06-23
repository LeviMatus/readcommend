import React from "react";
import * as models from "../../models";
import icon from "../../images/book.svg";
import { makeStyles } from "@material-ui/core/styles";
import { Rating } from "@material-ui/lab";
import { Chip } from "@material-ui/core";

const useStyles = makeStyles(theme => ({
  tile: {
    background: "#f6f6f6",
    marginTop: "8px",
    marginBottom: "8px",
    width: "100%",
    height: "90px",
    display: "flex",
    flexDirection: "row",
    flexWrap: "nowrap",
  },
  icon: {
    background: `url(${icon}) center no-repeat`,
    width: "80px",
    height: "80px",
    marginTop: "auto",
    marginBottom: "auto",
  },
  info: {
    padding: "8px",
  },
  title: {
    fontSize: "14pt",
  },
  author: {
    fontSize: "11pt",
    fontStyle: "italic",
    color: "#414141",
  },
  year: {
    fontSize: "10pt",
    fontStyle: "italic",
    color: "#888888",
    marginLeft: "0.5em",
  },
  genreRow: {
    marginTop: "3px",
  },
  pages: {
    fontSize: "10pt",
    color: "#a27e39",
    marginLeft: "0.5em",
  },
  ratingColumn: {
    marginLeft: "auto",
    marginTop: "auto",
    marginBottom: "auto",
    marginRight: "20px",
  },
  rating: {
    marginTop: "auto",
    marginBottom: "auto",
  },
}));

export default function Book(props: { book: models.Book }) {
  const classes = useStyles();
  const b = props.book;
  return (
    <div className={classes.tile}>
      <div className={classes.icon} />
      <div className={classes.info}>
        <div className={classes.title}>{b.title}</div>
        <div>
          <span className={classes.author}>
            {b.author.firstName} {b.author.lastName}
          </span>
          <span className={classes.year}>({b.yearPublished})</span>
        </div>
        <div className={classes.genreRow}>
          <Chip label={b.genre.title} size="small" />
          <span className={classes.pages}>{b.pages} pages</span>
        </div>
      </div>
      <div className={classes.ratingColumn}>
        <Rating
          className={classes.rating}
          value={b.rating}
          precision={0.1}
          readOnly
        />
      </div>
    </div>
  );
}
