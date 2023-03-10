import React, { useState, useEffect } from 'react';
import './Beatsort.css';

function Beatsort() {
    const [songs, setSongs] = useState([]);

    useEffect(() => {
        // Fetch the initial playlist-order.json file to get the current order of songs
        fetch('/playlist-order.json')
            .then(response => response.json())
            .then(data => setSongs(data.songs));
    }, []);

    const handleDrag = (event, index) => {
        // Update the order of songs when the user drags and drops a song
        event.preventDefault();
        const draggedSong = songs[index];
        setSongs(prevSongs => {
            const updatedSongs = [...prevSongs];
            updatedSongs.splice(index, 1);
            updatedSongs.splice(event.target.dataset.index, 0, draggedSong);
            return updatedSongs;
        });
    };

    const handleDragOver = event => {
        event.preventDefault();
    };

    const savePlaylistOrder = () => {
        // Save the current order of songs to the playlist-order.json file
        fetch('/save-playlist-order', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ songs })
        });
    };

    return (
        <div className="beatsort-container">
            <h1>Beatsort</h1>
            <p>Drag and drop songs to sort your playlist.</p>
            <ul className="song-list">
                {songs.map((song, index) => (
                    <li key={song.id} data-index={index} draggable onDrag={event => handleDrag(event, index)} onDragOver={handleDragOver}>
                        {song.title} - {song.artist}
                    </li>
                ))}
            </ul>
            <button onClick={savePlaylistOrder}>Save Playlist Order</button>
        </div>
    );
}

export default Beatsort;
