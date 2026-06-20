package com.dkprojects.rmi;

import java.rmi.RemoteException;
import java.rmi.server.UnicastRemoteObject;

public class MyServiceImpl extends UnicastRemoteObject implements MyService{

    private KVStore store;

    public MyServiceImpl() throws RemoteException {

    }

    public MyServiceImpl(KVStore store){
        this.store = store;
    }
    
    public int get(String key) {
        return store.get(key);
    }

    public void put(String key, String value) {
        store.put(key, value);
    }
}