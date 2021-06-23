import { FormControl, Grid, InputLabel } from "@material-ui/core";
import React, { useState, useEffect } from "react";
import Authors from "./Authors";
import Genres from "./Genres";
import Eras from "./Eras";
import Sizes from "./Sizes";
import Limit from "./Limit";
import * as models from "../../models";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
  criteriaBox: {
    background: "#f4ebd9",
    marginBottom: "24px",
    width: "100%",
    marginLeft: "0px",
    marginRight: "0px",
  },
}));

export default function Criteria(props: {
  onChange: (criteria: models.Criteria) => void;
}) {
  const classes = useStyles();
  const [authors, setAuthors] = useState<models.Author[]>([]);
  const [genres, setGenres] = useState<models.Genre[]>([]);
  const [size, setSize] = useState<models.Size>(models.DefaultSize);
  const [era, setEra] = useState<models.Era>(models.DefaultEra);
  const [limit, setLimit] = useState<number>(10);

  useEffect(() => {
    props.onChange({
      authors: authors.map(x => x.id),
      genres: genres.map(x => x.id),
      minYear: era.minYear,
      maxYear: era.maxYear,
      minPages: size.minPages,
      maxPages: size.maxPages,
      limit: limit,
    });
  }, [authors, genres, size, era, limit]);

  return (
    <Grid container spacing={2} className={classes.criteriaBox}>
      <Grid item xs={12}>
        <Authors onChange={x => setAuthors(x)} />
      </Grid>
      <Grid item xs={12}>
        <Genres onChange={x => setGenres(x)} />
      </Grid>
      <Grid item xs={4}>
        <FormControl fullWidth={true}>
          <InputLabel>Era</InputLabel>
          <Eras onChange={x => setEra(x)} />
        </FormControl>
      </Grid>
      <Grid item xs={4}>
        <FormControl fullWidth={true}>
          <InputLabel>Pages</InputLabel>
          <Sizes onChange={x => setSize(x)} />
        </FormControl>
      </Grid>
      <Grid item xs={4}>
        <FormControl fullWidth={true}>
          <InputLabel>Max results</InputLabel>
          <Limit onChange={x => setLimit(x)} />
        </FormControl>
      </Grid>
    </Grid>
  );
}
