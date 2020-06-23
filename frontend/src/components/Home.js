import React, { useState, useEffect } from "react";
import { Card } from 'react-bootstrap';
import { LinkContainer } from "react-router-bootstrap";
import { ListGroupItem, ListGroup } from 'react-bootstrap'
import { useAppContext } from "../libs/contextLibs";
import { onError } from "../libs/errorLib";
import { API } from "aws-amplify";
import "./Home.css";


export default function Home() {
  const [events, setEvents] = useState([]);
  const { isAuthenticated } = useAppContext();
  const [isLoading, setIsLoading] = useState(true);
  useEffect(() => {
    async function onLoad() {
      if (!isAuthenticated) {
        return;
      }
  
      try {
        const loadedEvents = await loadEvents();
        if (loadedEvents.events !== null) {
          setEvents(loadedEvents.events);
        }
      } catch (e) {
        onError(e);
      }
  
      setIsLoading(false);
    }
  
    onLoad();
  }, [isAuthenticated]);
  
  function loadEvents() {
    return API.get("events", "/events");
  }

  function renderEventsList(events) {
    return [{}].concat(events).map((event, i) =>
      i !== 0 ? (
        <LinkContainer key={event.event_id} to={`/events/${event.event_id}`}>
          <Card className="text-center">
            <Card.Header as="h2">
              {event.name}
            </Card.Header>
            <Card.Body>
              <p>Description: {event.description}</p>
              <p>Status: {event.status}</p>
              <p>Start Time: {event.schedule.start_time}</p>
              <p>End Time: {event.schedule.end_time}</p>
            </Card.Body>
          </Card>
        </LinkContainer>
      ) : (
        <LinkContainer key="new" to="/events/new">
          <ListGroupItem>
            <h4>
              <b>{"\uFF0B"}</b> Create a new event
            </h4>
          </ListGroupItem>
        </LinkContainer>
      )
    );
  }

  function renderLander() {
    return (
      <div className="lander">
        <h1>Events App</h1>
        <p>A simple note taking app</p>
      </div>
    );
  }

  function renderEvents() {
    return (
      <div className="notes">
        <h1>Your Events</h1>
        <p>If you want to change the contents of the event, click on the card corresponding to it</p>
        <ListGroup>
          {!isLoading && renderEventsList(events)}
        </ListGroup>
      </div>
    );
  }

  return (
    <div className="Home">
      {isAuthenticated ? renderEvents() : renderLander()}
    </div>
  );
}