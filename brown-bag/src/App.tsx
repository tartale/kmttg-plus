import React from 'react';
import './App.css';

function App() {
  return (
    <div className="App">
      <ShowListing shows={[
        { id: 1, title: 'Breaking Bad', time: '9:00 PM', channel: 'AMC', image: 'https://i.ytimg.com/vi/cRFC1O7IVy8/maxresdefault.jpg', description: 'A high school chemistry teacher diagnosed with inoperable lung cancer turns to manufacturing and selling methamphetamine in order to secure his family\'s future.' },
        { id: 2, title: 'Game of Thrones', time: '8:00 PM', channel: 'HBO', image: 'https://i.ytimg.com/vi/cRFC1O7IVy8/maxresdefault.jpg', description: 'Nine noble families fight for control over the lands of Westeros, while an ancient enemy returns after being dormant for millennia.' },
        { id: 3, title: 'The Sopranos', time: '7:00 PM', channel: 'HBO', image: 'https://i.ytimg.com/vi/cRFC1O7IVy8/maxresdefault.jpg', description: 'New Jersey mob boss Tony Soprano deals with personal and professional issues in his home and business life that affect his mental state, leading him to seek professional psychiatric counseling.' },
      ]} />
    </div>
  );
}

export default App;

const ShowListing = (props: { shows: any }) => {
  const {shows} = props;

  return (
    <div className="show-listing">
      {shows.map((show: any) => (
        <div className="show-item" key={show.id}>
          <div className="show-details">
            <div className="show-title">{show.title}</div>
            <div className="show-time">{show.time}</div>
            <div className="show-channel">{show.channel}</div>
            <div className="show-description">{show.description}</div>
          </div>
        </div>
      ))}
    </div>
  );
};


