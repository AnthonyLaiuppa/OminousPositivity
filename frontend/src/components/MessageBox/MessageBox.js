import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Spinner } from 'react-bootstrap';

function MessageBox() {
    const [message, setMessage] = useState('');
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        fetch(`${process.env.REACT_APP_API_URL}/message`)
            .then(response => response.json())
            .then(data => {
                setMessage(data.message); // Adjust according to the API response structure
                setIsLoading(false);
            })
            .catch(error => {
                console.error('Error fetching data: ', error);
                setIsLoading(false);
            });
    }, []);

    return (
        <Container className="full-height">
            <Row>
                <Col>
                    {isLoading ? (
                        <Spinner animation="border" role="status">
                            <span className="visually-hidden">Loading...</span>
                        </Spinner>
                    ) : (
                        <div className="quote-box">
                            <p className="quote-text">"{message}"</p>
                        </div>
                    )}
                </Col>
            </Row>
        </Container>
    );
}

export default MessageBox;
