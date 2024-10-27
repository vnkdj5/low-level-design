package elevator;

public class Request implements Comparable<Request> {
	private int sourceFloor;
    private int destinationFloor;
	private Direction direction;

    public Request(int sourceFloor, int destinationFloor) {
        this.sourceFloor = sourceFloor;
        this.destinationFloor = destinationFloor;
		if(sourceFloor<destinationFloor){
			this.direction = Direction.UP;
		}else{
			this.direction = Direction.DOWN;
		}
    }

    public int getDestinationFloor() {
        return destinationFloor;
    }

    public int getSourceFloor() {
        return sourceFloor;
    }

	public Direction getDirectionToGo(){
		return this.direction;
	}


	@Override
    public int compareTo(Request req) {
        return Integer.compare(this.destinationFloor, req.destinationFloor);
    }
   /*  private InternalRequest internalRequest;
	private ExternalRequest externalRequest;
	public Request(InternalRequest internalRequest, ExternalRequest externalRequest) {
		this.internalRequest = internalRequest;
		this.externalRequest = externalRequest;
	}

    public InternalRequest getInternalRequest() {
		return internalRequest;
	}
	public void setInternalRequest(InternalRequest internalRequest) {
		this.internalRequest = internalRequest;
	}
	public ExternalRequest getExternalRequest() {
		return externalRequest;
	}
	public void setExternalRequest(ExternalRequest externalRequest) {
		this.externalRequest = externalRequest;
	}

    public int getDestinationFloor(){
        return this.getInternalRequest().getDestinationFloor();
    }

    public int getSourceFloor(){
        return this.getExternalRequest().getSourceFloor();
    }

	@Override
	public int compareTo(Request req) {
		if (this.getInternalRequest().getDestinationFloor() == req.getInternalRequest().getDestinationFloor())
			return 0;
		else if (this.getInternalRequest().getDestinationFloor() > req.getInternalRequest().getDestinationFloor())
			return 1;
		else
			return -1;
	} */
}
