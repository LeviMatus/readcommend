import React, { useEffect, useState } from "react";
import * as models from "../../models";
import Book from "./Book";

export default function Results(props: { criteria: models.Criteria }) {
  const [error, setError] = useState<models.Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [books, setBooks] = React.useState<models.Book[]>([]);

  useEffect(() => {
    const c = props.criteria;
    const url = getUrl(
      "http://localhost:5000/api/v1/books",
      getListParam("authors", c.authors),
      getListParam("genres", c.genres),
      getParam("min-pages", c.minPages),
      getParam("max-pages", c.maxPages),
      getParam("min-year", c.minYear),
      getParam("max-year", c.maxYear),
      getParam("limit", c.limit)
    );

    console.log("Fetching:", url);
    fetch(url)
      .then(result => {
        result.json().then(body => {
          if (result.ok) {
            setIsLoaded(true);
            setBooks(body);
          } else {
            setIsLoaded(true);
            setError(body);
          }
        });
      })
      .catch(error => {
        setIsLoaded(true);
        setError(error);
      });
  }, [props.criteria]);

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  if (!isLoaded) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      {books.map(x => (
        <Book key={x.id} book={x} />
      ))}
    </div>
  );
}

function getUrl(path: string, ...values: (string | undefined)[]): string {
  return path + "?" + values.filter(x => x !== undefined).join("&");
}

function getListParam(name: string, value: any[]): string | undefined {
  return value.length > 0 ? `${name}=${value.join(",")}` : undefined;
}

function getParam(name: string, value?: any): string | undefined {
  return value !== undefined && value !== null ? `${name}=${value}` : undefined;
}
