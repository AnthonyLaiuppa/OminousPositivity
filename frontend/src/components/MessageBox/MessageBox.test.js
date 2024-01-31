import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import MessageBox from './MessageBox';

// Mock the global fetch API before all tests
beforeAll(() => {
    global.fetch = jest.fn();
});

// Reset all mocks after each test
afterEach(() => {
    jest.resetAllMocks();
});

describe('<MessageBox />', () => {
    test('displays the loading spinner before the message is fetched', async () => {
        // Mock the fetch response
        global.fetch.mockResolvedValueOnce({
            json: () => Promise.resolve({ message: 'Test message' }),
        });

        render(<MessageBox />);

        // Check if the spinner is in the document
        expect(screen.getByRole('status')).toBeInTheDocument();

        // Wait for and check if the message is eventually displayed in the document
        await waitFor(() => expect(screen.getByText('"Test message"')).toBeInTheDocument());
    });

    test('fetches and displays a message', async () => {
        // Mock the fetch response
        global.fetch.mockResolvedValueOnce({
            json: () => Promise.resolve({ message: 'Test message' }),
        });

        render(<MessageBox />);

        // Wait for and check if the message is eventually displayed in the document
        await waitFor(() => expect(screen.getByText('"Test message"')).toBeInTheDocument());
    });

});
