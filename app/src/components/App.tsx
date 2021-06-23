import React from "react";
import Criteria from "./criteria/Criteria";
import Results from "./results/Results";
import Container from "@material-ui/core/Container";
import Box from "@material-ui/core/Box";
import * as models from "../models";
import Copyright from "./Copyright";

export default function App() {
  const [criteria, setCriteria] = React.useState<models.Criteria>(
    models.DefaultCriteria
  );

  return (
    <Container maxWidth="md">
      <Box my={4}>
        <Criteria onChange={x => setCriteria(x)} />
        <Results criteria={criteria} />
        <Copyright />
      </Box>
    </Container>
  );
}
