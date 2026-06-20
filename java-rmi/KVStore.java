package com.dkprojects.rmi;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class KVStore{
    private Map<String, String> store = new ConcurrentHashMap();
    
    public int get(String key) {
        return store.get(key);
    }

    public void put(String key, String value) {
        store.put(key, value);
    }
}