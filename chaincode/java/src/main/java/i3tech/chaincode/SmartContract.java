package i3tech.chaincode;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.hyperledger.fabric.shim.ChaincodeBase;
import org.hyperledger.fabric.shim.ChaincodeStub;
import org.hyperledger.fabric.shim.ledger.KeyModification;
import org.hyperledger.fabric.shim.ledger.QueryResultsIterator;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.Date;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.IntStream;

public class SmartContract extends ChaincodeBase {


    private static Log LOG = LogFactory.getLog(SmartContract.class);

    public static final String INVOKE_FUNCTION = "invoke";
    public static final String QUERY_FUNCTION = "query";
    public static final String QUERY_HISTORY_FUNCTION = "queryHistory";

    @Override
    public Response init(ChaincodeStub chaincodeStub) {
        return newSuccessResponse();
    }

    @Override
    public Response invoke(ChaincodeStub chaincodeStub) {

        String functionName = chaincodeStub.getFunction();
        LOG.info("function name: "+ functionName);


        List<String> paramList = chaincodeStub.getParameters();
        IntStream.range(0,paramList.size()).forEach(idx -> LOG.info("value of param: " + idx  + " is: "+paramList.get(idx)));

        if (INVOKE_FUNCTION.equalsIgnoreCase(functionName)) {
            return invokeOperation(chaincodeStub, paramList);
        } else if (QUERY_FUNCTION.equalsIgnoreCase(functionName)) {
            return queryOperation(chaincodeStub, paramList);
        }   else if (QUERY_HISTORY_FUNCTION.equalsIgnoreCase(functionName)){
                return queryByHistoryFunction(chaincodeStub, paramList);
            }
         else return newErrorResponse(functionName + " function is currently not supported");
    }

    private Response queryByHistoryFunction(ChaincodeStub chaincodeStub, List<String> paramList) {
        QueryResultsIterator<KeyModification> queryResultsIterator = chaincodeStub.getHistoryForKey(paramList.get(0));
        return newSuccessResponse(buildJsonFromQueryResult(queryResultsIterator));

    }

    private String buildJsonFromQueryResult(QueryResultsIterator<KeyModification> queryResultsIterator) {

        JSONArray jsonArray = new JSONArray();
        queryResultsIterator.forEach(keyModification -> {
            Map<String, Object> map = new LinkedHashMap<>();
            map.put("transactionId", keyModification.getTxId());
            map.put("timestamp", keyModification. getTimestamp().toString());
            map.put("value", keyModification.getStringValue());
            map.put("isDeleted", keyModification.isDeleted());
            jsonArray.put(map);
        });

        JSONObject jsonObject = new JSONObject();
        try {
            jsonObject.accumulate("transactions", jsonArray);
        } catch (JSONException e) {
            throw new RuntimeException("exception while generating json object");
        }
        return jsonObject.toString();
    }

    private Response queryOperation(ChaincodeStub chaincodeStub, List<String> paramList) {

        String misteryBox = chaincodeStub.getStringState(paramList.get(0));
        if (Objects.isNull(misteryBox)) {
            return newErrorResponse("mileage of provided car not found");
        }
        return newSuccessResponse(misteryBox);
    }

    private Response invokeOperation(ChaincodeStub chaincodeStub, List<String> paramList) {

        Misterybox mbox = new Misterybox();
        mbox.setSerial(paramList.get(0));
        mbox.setSize(paramList.get(1));
        mbox.setOwner(paramList.get(2));


        String misteryBoxFromLedger = chaincodeStub.getStringState(mbox.getSerial());

        if (StringUtils.isEmpty(misteryBoxFromLedger)) {
            chaincodeStub.putStringState(mbox.getSerial(), mbox.toJSON());
        } else {
            if (Integer.valueOf(misteryBoxFromLedger).compareTo(Integer.valueOf(mbox.toJSON())) >= 0) {
                return newErrorResponse("incorrect value");
            }
            chaincodeStub.putStringState(mbox.getSerial(), mbox.toJSON());
        }
        return newSuccessResponse();
    }


    public static void main(String [] args){
        new SmartContract().start(args);
    }

}

class Misterybox {
    private String docType = "MisteryBox";
    private String serial;
    private String size;
    private String owner;
    private Date registerAt;


    public String getDocType() {
        return this.docType;
    }

    public String getSerial() {
        return this.serial;
    }

    public void setSerial(String serial) {
        this.serial = serial;
    }

    public String getSize() {
        return this.size;
    }

    public void setSize(String size) {
        this.size = size;
    }

    public String getOwner() {
        return this.owner;
    }

    public void setOwner(String owner) {
        this.owner = owner;
    }

    public Date getRegisterAt() {
        return this.registerAt;
    }

    public void setRegisterAt(Date registerAt) {
        this.registerAt = registerAt;
    }

    public String toJSON(){
        return new Gson().toJson(this);
    }

    public Misterybox(){
        this.registerAt = new Date();
    }


}
