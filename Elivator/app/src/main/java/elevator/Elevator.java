package elevator;

import java.util.TreeSet;

public class Elevator implements Runnable {

    private int currentFloor = 0;
    private Direction currentDirection = Direction.UP;
    private State currentState = State.IDLE;

    /**
     * jobs which are being processed
     */
    private TreeSet<Request> currentJobs = new TreeSet<Request>();

    /**
     * up jobs which cannot be processed now so put in pending queue
     */
    private TreeSet<Request> upPendingJobs = new TreeSet<Request>();

    /**
     * down jobs which cannot be processed now so put in pending queue
     */
    private TreeSet<Request> downPendingJobs = new TreeSet<Request>();

    public void startElevator() {
        System.out.println("The Elevator has started functioning");
        while (true) {
            Request currentRequest = null; // Store the request being processed
            synchronized (this) { // Acquire lock for accessing shared state
                if (checkIfJob()) {
                    System.out.println("There are pending jobs");
                    if (currentDirection == Direction.UP) {
                        currentRequest = currentJobs.pollFirst();
                    } else {
                        currentRequest = currentJobs.pollLast();
                    }
                } else {
                    try {
                        System.out.println("Waiting for jobs...");
                        wait(); // Release lock and wait for a new job
                    } catch (InterruptedException e) {
                        System.err.println("Elevator thread interrupted.");
                        break; // Exit the loop on interruption
                    }
                }
            } // End of synchronized block

            // Process the job outside the synchronized block
            if (currentRequest != null) {
                if (currentDirection == Direction.UP) {
                    processUpRequest(currentRequest);
                    addPendingDownJobsToCurrentJobs();

                    // ... (rest of your processing logic) ...
                } else {
                    processDownRequest(currentRequest);
                    addPendingUpJobsToCurrentJobs();

                    // ... (rest of your processing logic) ...
                }
            }
        }
        System.out.println("Elevator stopped.");
    }

    // Check if there are pending jobs
    public boolean checkIfJob() {
        //System.out.println(currentJobs.size());
        if (currentJobs.isEmpty()) {
            return false;
        }
       

        return true;

    }

    private void processUpRequest(Request request) {
        // The elevator is not on the floor where the person has requested it i.e.
        // source floor. So first bring it there.
        int start = currentFloor;
        if (start < request.getSourceFloor()) {
            for (int i = start; i < request.getSourceFloor(); i++) {
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    // TODO Auto-generated catch block
                    e.printStackTrace();
                }
                System.out.println("We have reached floor -- " + i);
                currentFloor = i;
            }

        }

        // The elevator is now on the floor where the person has requested it i.e.
        // source floor. User can enter and go to the destination floor.
        System.out.println("Reached Source Floor--opening door");

        start = currentFloor;
        // Go to destination floor
        for (int i = start + 1; i <= request.getDestinationFloor(); i++) {
            try {
                Thread.sleep(1000);
            } catch (InterruptedException e) {
                // TODO Auto-generated catch block
                e.printStackTrace();
            }
            System.out.println("We have reached floor -- " + i);
            currentFloor = i;
            // check if intermediate request can be processed
            if (checkIfNewJobCanBeProcessed(request)) {
                break;
            }
        }

    }

    private void processDownRequest(Request request) {
        // The elevator is not on the floor where the person has requested it i.e.
        // source floor. So first bring it there.
        int start = currentFloor;

        if (start < request.getSourceFloor()) {
            for (int i = start; i < request.getSourceFloor(); i++) {
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    // TODO Auto-generated catch block
                    e.printStackTrace();
                }
                System.out.println("We have reached floor -- " + i);
                currentFloor = i;
            }
        }
        System.out.println("Reached Source Floor--opening door");

        start = currentFloor;

        for (int i = start - 1; i >= request.getDestinationFloor(); i--) {
            try {
                Thread.sleep(1000);
            } catch (InterruptedException e) {
                // TODO Auto-generated catch block
                e.printStackTrace();
            }
            System.out.println("We have reached floor -- " + i);
            currentFloor = i;

            // check if intermediate request can be processed
            if (checkIfNewJobCanBeProcessed(request)) {
                break;
            }

        }

    }

    private boolean checkIfNewJobCanBeProcessed(Request currentRequest) {
        if (checkIfJob()) {
            if (currentDirection == Direction.UP) {

                Request request = currentJobs.pollFirst(); // lowest floor
                if (request.getDestinationFloor() < currentRequest.getDestinationFloor()) {
                    currentJobs.add(request);
                    currentJobs.add(currentRequest);
                    return true;
                } else {
                    currentJobs.add(request);
                }
            }

            if (currentDirection == Direction.DOWN) {

                Request request = currentJobs.pollLast(); // highest floor

                if (request.getDestinationFloor() > currentRequest.getDestinationFloor()) {
                    currentJobs.add(request);
                    currentJobs.add(currentRequest);
                    return true;
                } else {
                    currentJobs.add(request);
                }

            }
        }

        return false;

    }

    private void addPendingDownJobsToCurrentJobs() {
        if (!downPendingJobs.isEmpty()) {
            System.out.println("Pick a pending down job and execute it by putting in current job");

            currentJobs = downPendingJobs;
            currentDirection = Direction.DOWN;
        } else {
            currentState = State.IDLE;
            System.out.println("addPendingDownJobsToCurrentJobs: The elevator is in Idle state");

        }

    }

    private void addPendingUpJobsToCurrentJobs() {
        if (!upPendingJobs.isEmpty()) {
            System.out.println("Pick a pending up job and execute it by putting in current job");

            currentJobs = upPendingJobs;
            currentDirection = Direction.UP;
        } else {
            currentState = State.IDLE;
            System.out.println("addPendingUpJobsToCurrentJobs: The elevator is in Idle state");

        }

    }

    public synchronized void addJob(Request request) {
        System.err.println("REquest Dest:"+request.getDestinationFloor()+" src:"+request.getSourceFloor());
        if (currentState == State.IDLE) {
            System.out.println("Changing state to moving");
            currentState = State.MOVING;
            currentDirection = request.getDirectionToGo();
            currentJobs.add(request);
        } else if (currentState == State.MOVING) {
            System.out.println("currentState:  moving");
            if (request.getDirectionToGo() != currentDirection) {
                addtoPendingJobs(request);
            } else {
                if (currentDirection == Direction.UP && request.getDestinationFloor() < currentFloor) {
                    addtoPendingJobs(request);
                    return;
                } else if (currentDirection == Direction.DOWN && request.getDestinationFloor() > currentFloor) {
                    addtoPendingJobs(request);

                } else {
                    currentJobs.add(request);

                }

            }

        }
        //startElevator();
        notifyAll();

    }

    public void addtoPendingJobs(Request request) {
        if (request.getDirectionToGo() == Direction.UP) {
            System.out.println("Add to pending up jobs");
            upPendingJobs.add(request);
        } else {
            System.out.println("Add to pending down jobs");
            downPendingJobs.add(request);
        }
    }

    public void run(){
        startElevator();
    }

}
