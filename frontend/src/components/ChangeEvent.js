import React, { useState } from "react";
import { FormGroup, FormControl, FormLabel } from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import { onError } from "../libs/errorLib";
import { useParams, useHistory } from "react-router-dom";
// import config from "../config";
import "./NewEvent.css";
import { API } from "aws-amplify";

export default function ChangeEvent() {
  const history = useHistory();
  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [status, setStatus] = useState("")
  const [startTime, setStartTime] = useState("")
  const [endTime, setEndTime] = useState("")
  const { event_id } = useParams()
  const [isLoading, setIsLoading] = useState(false);

  function validateForm() {
    return name.length > 0 || description.length > 0 || status.length > 0 || startTime.length > 0 || endTime.length > 0;
  }

  function patchEvent(eventObj) {
    return API.patch("events", `/events/${event_id}`, {
      body: eventObj
    })
  }

  async function handleSubmit(event) {
      event.preventDefault();
      setIsLoading(true);
      try {
        const eventObj = {};
        if(name.length > 0) {
          eventObj["name"] = name;
        }
        if(description.length > 0) {
          eventObj["description"] = description;
        }
        if(status.length > 0) {
          eventObj["status"] = status;
        }
        if(startTime.length > 0) {
          if (!eventObj.hasOwnProperty("schedule")) {
            eventObj["schedule"] = {}
          }
          eventObj["schedule"].start_time = startTime;
        }
        if(endTime.length > 0) {
          if (!eventObj.hasOwnProperty("schedule")) {
            eventObj["schedule"] = {}
          }
          eventObj["schedule"].end_time = endTime;
        }
        await patchEvent(eventObj)
        history.push("/");
      } catch(e) {
        onError(e);
        setIsLoading(false);
      }
  }

  return (
    <div className="NewEvent">
      <form onSubmit={handleSubmit}>
        <FormLabel>Name</FormLabel>
        <FormGroup controlId="name">
          <FormControl 
            value={name}
            componentclass="text"
            onChange={e => setName(e.target.value)}
          />
        </FormGroup>
        <FormGroup controlId="description">
          <FormLabel>Description</FormLabel>
          <FormControl 
            value={description}
            componentclass="text"
            onChange={e => setDescription(e.target.value)}
          />
        </FormGroup>
        <FormGroup controlId="status">
          <FormLabel>Status</FormLabel>
          <FormControl 
            value={status}
            componentclass="text"
            onChange={e => setStatus(e.target.value)}
          />
        </FormGroup>
        <FormGroup controlId="startTime">
          <FormLabel>Start Time</FormLabel>
          <FormControl 
            value={startTime}
            componentclass="text"
            onChange={e => setStartTime(e.target.value)}
          />
        </FormGroup>
        <FormGroup controlId="endTime">
          <FormLabel>End Time</FormLabel>
          <FormControl 
            value={endTime}
            componentclass="textarea"
            onChange={e => setEndTime(e.target.value)}
          />
        </FormGroup>
        <LoaderButton
          block
          type="submit"
          bssize="large"
          bsstyle="primary"
          isLoading={isLoading}
          disabled={!validateForm()}
        >
          Change the event
        </LoaderButton>
      </form>
    </div>
  );
}